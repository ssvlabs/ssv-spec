package frost

import (
	"math/rand"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	ecies "github.com/ecies/go/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

var (
	thisCurve  = curves.BLS12381G1()
	domainType = types.PrimusTestnet
	sigType    = types.DKGSignatureType
)

func init() {
	types.InitBLS()
}

// FROST contains network to broadcast messages, signer to sign outgoing messages,
// storage to store keygen output, config to store protocol configuration such as id,
// threshold, operator list etc and state contains properties that changes over the
// runtime of the protocol like current round, msgs etc
type FROST struct {
	network dkg.Network
	signer  types.DKGSigner
	storage dkg.Storage
	config  ProtocolConfig
	state   *ProtocolState
}

// ProtocolConfig contains properties needed to start the protocol like requestID,
// operatorID, threshold, operator list etc.
type ProtocolConfig struct {
	identifier      dkg.RequestID
	threshold       uint32
	operatorID      types.OperatorID
	operators       []uint32
	operatorsOld    []uint32
	oldKeyGenOutput *dkg.KeyGenOutput
}

func (c *ProtocolConfig) isResharing() bool {
	return len(c.operatorsOld) > 0
}

func (c *ProtocolConfig) inOldCommittee() bool {
	for _, id := range c.operatorsOld {
		if types.OperatorID(id) == c.operatorID {
			return true
		}
	}
	return false
}

func (c *ProtocolConfig) inNewCommittee() bool {
	for _, id := range c.operators {
		if types.OperatorID(id) == c.operatorID {
			return true
		}
	}
	return false
}

// ProtocolState maintains value for current round, messages, sessions key and
// operator shares. these properties will change over the runtime of the protocol
type ProtocolState struct {
	currentRound   ProtocolRound
	participant    *frost.DkgParticipant
	sessionSK      *ecies.PrivateKey
	msgContainer   *MsgContainer
	operatorShares map[uint32]*bls.SecretKey
}

// ProtocolRound is enum for all the rounds in the protocol
type ProtocolRound int

const (
	Uninitialized ProtocolRound = iota
	Preparation
	Round1
	Round2
	KeygenOutput
	Blame
)

var rounds = []ProtocolRound{
	Uninitialized,
	Preparation,
	Round1,
	Round2,
	KeygenOutput,
	Blame,
}

// New creates a new protocol instance for new keygen
// TODO: func New(network dkg.Network, signer types.DKGSigner, storage dkg.Storage, config ProtocolConfig) {}
func New(
	network dkg.Network,
	operatorID types.OperatorID,
	requestID dkg.RequestID,
	signer types.DKGSigner,
	storage dkg.Storage,
	init *dkg.Init,
) dkg.Protocol {

	config := ProtocolConfig{
		identifier: requestID,
		threshold:  uint32(init.Threshold),
		operatorID: operatorID,
		operators:  types.OperatorList(init.OperatorIDs).ToUint32List(),
	}
	return newProtocol(network, signer, storage, config)
}

// NewResharing creates a new protocol instance for resharing
// Temporary, TODO: Remove and use interface with Reshare
func NewResharing(
	network dkg.Network,
	operatorID types.OperatorID,
	requestID dkg.RequestID,
	signer types.DKGSigner,
	storage dkg.Storage,
	operatorsOld []types.OperatorID,
	init *dkg.Reshare,
	output *dkg.KeyGenOutput,
) dkg.Protocol {

	config := ProtocolConfig{
		identifier:      requestID,
		threshold:       uint32(init.Threshold),
		operatorID:      operatorID,
		operators:       types.OperatorList(init.OperatorIDs).ToUint32List(),
		operatorsOld:    types.OperatorList(operatorsOld).ToUint32List(),
		oldKeyGenOutput: output,
	}
	return newProtocol(network, signer, storage, config)
}

func newProtocol(network dkg.Network, signer types.DKGSigner, storage dkg.Storage, config ProtocolConfig) dkg.Protocol {
	return &FROST{
		network: network,
		signer:  signer,
		storage: storage,
		config:  config,
		state: &ProtocolState{
			currentRound:   Uninitialized,
			msgContainer:   newMsgContainer(),
			operatorShares: make(map[uint32]*bls.SecretKey),
		},
	}
}

// Start initializes frost participant, generates a session key pair and broadcasts preparation message.
// TODO: If Reshare, confirm participating operators using qbft before kick-starting this process.
func (fr *FROST) Start() error {
	fr.state.currentRound = Preparation

	// create a new dkg participant
	ctx := make([]byte, 16)
	if _, err := rand.Read(ctx); err != nil {
		return err
	}
	participant, err := frost.NewDkgParticipant(uint32(fr.config.operatorID), fr.config.threshold, string(ctx), thisCurve, fr.config.operators...)
	if err != nil {
		return errors.Wrap(err, "failed to initialize a dkg participant")
	}
	fr.state.participant = participant

	if !fr.needToRunCurrentRound() {
		return nil // preparation round only runs in new committee
	}

	// generate session key pair
	k, err := ecies.GenerateKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate session sk")
	}
	fr.state.sessionSK = k

	// create and broadcast PreparationMessage
	msg := &ProtocolMsg{
		Round: Preparation,
		PreparationMessage: &PreparationMessage{
			SessionPk: k.PublicKey.Bytes(true),
		},
	}
	_, err = fr.broadcastDKGMessage(msg)
	return err
}

// ProcessMsg  decodes and validates incoming message. It then check for blame
// or proceed with processing these messages based on their round.
func (fr *FROST) ProcessMsg(msg *dkg.SignedMessage) (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {

	// validated signed message
	if err := fr.validateSignedMessage(msg); err != nil {
		return false, nil, errors.Wrap(err, "failed to Validate signed message")
	}

	// decodes and validates protocol message
	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "failed to decode protocol msg")
	}
	if err := protocolMessage.Validate(); err != nil {
		return fr.createAndBroadcastBlameOfInvalidMessage(uint32(msg.Signer), msg)
	}

	// store incoming message unless it already exists, then check for blame
	existingMessage, err := fr.state.msgContainer.SaveMsg(protocolMessage.Round, msg)
	if err != nil && !haveSameRoot(existingMessage, msg) {
		return fr.createAndBroadcastBlameOfInconsistentMessage(existingMessage, msg)
	}

	// process message based on their round
	switch protocolMessage.Round {
	case Preparation:
		return fr.processRound1()
	case Round1:
		return fr.processRound2()
	case Round2:
		return fr.processKeygenOutput()
	case Blame:
		// here we are checking blame right away unlike other rounds where
		// we wait to receive messages from all the operators in the protocol
		return fr.checkBlame(uint32(msg.Signer), protocolMessage, msg)
	default:
		return true, nil, dkg.ErrInvalidRound{}
	}
}

func (fr *FROST) canProceedThisRound() bool {
	// Note: for Resharing, Preparation (New Committee) -> Round1 (Old Committee) -> Round2 (New Committee)
	if fr.config.isResharing() && fr.state.currentRound == Round1 {
		return fr.state.msgContainer.allMessagesReceivedFor(Round1, fr.config.operatorsOld)
	}
	return fr.state.msgContainer.allMessagesReceivedFor(fr.state.currentRound, fr.config.operators)
}

func (fr *FROST) needToRunCurrentRound() bool {
	if !fr.config.isResharing() {
		return true // always run for new keygen
	}
	switch fr.state.currentRound {
	case Preparation, Round2, KeygenOutput:
		return fr.config.inNewCommittee()
	case Round1:
		return fr.config.inOldCommittee()
	default:
		return false
	}
}

func (fr *FROST) validateSignedMessage(msg *dkg.SignedMessage) error {
	if msg.Message.Identifier != fr.config.identifier {
		return errors.New("got mismatching identifier")
	}

	found, operator, err := fr.storage.GetDKGOperator(msg.Signer)
	if !found {
		return errors.New("unable to find signer")
	}
	if err != nil {
		return errors.Wrap(err, "unable to find signer")
	}

	root, err := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domainType, sigType))
	if err != nil {
		return errors.Wrap(err, "failed to get root")
	}

	pk, err := crypto.Ecrecover(root, msg.Signature)
	if err != nil {
		return errors.Wrap(err, "unable to recover public key")
	}

	addr := common.BytesToAddress(crypto.Keccak256(pk[1:])[12:])
	if addr != operator.ETHAddress {
		return errors.New("invalid signature")
	}
	return nil
}
