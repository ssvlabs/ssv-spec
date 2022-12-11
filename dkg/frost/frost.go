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

type FROST struct {
	network dkg.Network
	signer  types.DKGSigner
	storage dkg.Storage
	config  ProtocolConfig
	state   *ProtocolState
}

type ProtocolState struct {
	currentRound   ProtocolRound
	participant    *frost.DkgParticipant
	msgs           ProtocolMessageStore
	operatorShares map[uint32]*bls.SecretKey
}

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

func New(
	network dkg.Network,
	operatorID types.OperatorID,
	requestID dkg.RequestID,
	signer types.DKGSigner,
	storage dkg.Storage,
	init *dkg.Init,
) dkg.Protocol {

	fr := &FROST{
		network: network,
		signer:  signer,
		storage: storage,
		config: ProtocolConfig{
			identifier: requestID,
			threshold:  uint32(init.Threshold),
			operatorID: operatorID,
			operators:  types.OperatorList(init.OperatorIDs).ToUint32List(),
		},
		state: &ProtocolState{
			currentRound:   Uninitialized,
			msgs:           newProtocolMessageStore(),
			operatorShares: make(map[uint32]*bls.SecretKey),
		},
	}

	return fr
}

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

	return &FROST{
		network: network,
		signer:  signer,
		storage: storage,
		config: ProtocolConfig{
			identifier:      requestID,
			threshold:       uint32(init.Threshold),
			operatorID:      operatorID,
			operators:       types.OperatorList(init.OperatorIDs).ToUint32List(),
			operatorsOld:    types.OperatorList(operatorsOld).ToUint32List(),
			oldKeyGenOutput: output,
		},
		state: &ProtocolState{
			currentRound:   Uninitialized,
			msgs:           newProtocolMessageStore(),
			operatorShares: make(map[uint32]*bls.SecretKey),
		},
	}
}

// TODO: If Reshare, confirm participating operators using qbft before kick-starting this process.
func (fr *FROST) Start() error {
	fr.state.currentRound = Preparation

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
		return nil
	}

	k, err := ecies.GenerateKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate session sk")
	}
	fr.config.sessionSK = k

	msg := &ProtocolMsg{
		Round: fr.state.currentRound,
		PreparationMessage: &PreparationMessage{
			SessionPk: k.PublicKey.Bytes(true),
		},
	}

	_, err = fr.broadcastDKGMessage(msg)
	return err
}

func (fr *FROST) ProcessMsg(msg *dkg.SignedMessage) (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {

	if err := fr.validateSignedMessage(msg); err != nil {
		return false, nil, errors.Wrap(err, "failed to Validate signed message")
	}

	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "failed to decode protocol msg")
	}

	if err := protocolMessage.Validate(); err != nil {
		return fr.createAndBroadcastBlameOfInvalidMessage(uint32(msg.Signer), msg)
	}

	existingMessage, ok := fr.state.msgs[protocolMessage.Round][uint32(msg.Signer)]
	isBlameTypeInconsisstent := (ok && !fr.haveSameRoot(existingMessage, msg))
	if isBlameTypeInconsisstent {
		return fr.createAndBroadcastBlameOfInconsistentMessage(existingMessage, msg)
	}

	fr.state.msgs[protocolMessage.Round][uint32(msg.Signer)] = msg

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
		return fr.state.msgs.allMessagesReceivedFor(Round1, fr.config.operatorsOld)
	}
	return fr.state.msgs.allMessagesReceivedFor(fr.state.currentRound, fr.config.operators)
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
