package frost

import (
	"math/rand"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	ecies "github.com/ecies/go/v2"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

// Instance is a single FROST DKG instance. It implements Start to init DKG
// parameters and broadcast session key, and ProcessMsg to process 2 rounds of
// frost, blame and keygen output messages
type Instance struct {
	state  *State
	config dkg.IConfig

	instanceParams InstanceParams
}

// New creates a new protocol instance for new keygen
// TODO: func New(network dkg.Network, signer types.DKGSigner, storage dkg.Storage, config ProtocolConfig) {}
func New(
	requestID dkg.RequestID,
	operatorID types.OperatorID,
	config dkg.IConfig,
	init *dkg.Init,
) dkg.Protocol {

	instanceParams := InstanceParams{
		identifier: requestID,
		threshold:  uint32(init.Threshold),
		operatorID: operatorID,
		operators:  types.OperatorList(init.OperatorIDs).ToUint32List(),
	}
	return newProtocol(config, instanceParams)
}

// NewResharing creates a new protocol instance for resharing
// Temporary, TODO: Remove and use interface with Reshare
func NewResharing(
	requestID dkg.RequestID,
	operatorID types.OperatorID,
	config dkg.IConfig,
	reshare *dkg.Reshare,
	reshareParams *dkg.ReshareParams,
) dkg.Protocol {

	instanceParams := InstanceParams{
		identifier:      requestID,
		threshold:       uint32(reshare.Threshold),
		operatorID:      operatorID,
		operators:       types.OperatorList(reshare.OperatorIDs).ToUint32List(),
		operatorsOld:    types.OperatorList(reshare.OldOperatorIDs).ToUint32List(),
		oldKeyGenOutput: reshareParams.OldKeygenOutput,
	}
	return newProtocol(config, instanceParams)
}

func newProtocol(config dkg.IConfig, instanceParams InstanceParams) dkg.Protocol {
	return &Instance{
		config:         config,
		instanceParams: instanceParams,

		state: initState(),
	}
}

// Start initializes frost participant, generates a session key pair and
// broadcasts preparation message.
// TODO: If Reshare, confirm participating operators using qbft before kick-
// starting this process.
func (fr *Instance) Start() error {
	fr.state.roundTimer.OnTimeout(fr.UponRoundTimeout)
	fr.state.SetCurrentRound(common.Preparation)
	fr.state.roundTimer.StartRoundTimeoutTimer(fr.state.GetCurrentRound())

	// create a new dkg participant
	ctx := make([]byte, 16)
	if _, err := rand.Read(ctx); err != nil {
		return err
	}
	participant, err := frost.NewDkgParticipant(uint32(fr.instanceParams.operatorID), fr.instanceParams.threshold, string(ctx), thisCurve, fr.instanceParams.operators...)
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
		Round: common.Preparation,
		PreparationMessage: &PreparationMessage{
			SessionPk: k.PublicKey.Bytes(true),
		},
	}
	bcastMsg, err := fr.saveSignedMsg(msg)
	if err != nil {
		return err
	}
	return fr.config.GetNetwork().BroadcastDKGMessage(bcastMsg)
}

// ProcessMsg  decodes and validates incoming message. It then check for blame
// to create and broadcast blame message or proceeds with processing 2 rounds of
// frost, blame or keygen-output as per the message's round
func (fr *Instance) ProcessMsg(msg *dkg.SignedMessage) (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {

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
	case common.Preparation:
		return fr.processRound1()
	case common.Round1:
		return fr.processRound2()
	case common.Round2:
		return fr.processKeygenOutput()
	case common.Blame:
		// here we are checking blame right away unlike other rounds where
		// we wait to receive messages from all the operators in the protocol
		return fr.checkBlame(uint32(msg.Signer), protocolMessage, msg)
	case common.Timeout:
		return fr.ProcessTimeoutMessage()
	default:
		return true, nil, dkg.ErrInvalidRound{}
	}
}

func (fr *Instance) canProceedThisRound() bool {
	// Note: for Resharing, Preparation (New Committee) -> Round1 (Old Committee) -> Round2 (New Committee)
	if fr.instanceParams.isResharing() && fr.state.GetCurrentRound() == common.Round1 {
		return fr.state.msgContainer.AllMessagesReceivedFor(common.Round1, fr.instanceParams.operatorsOld)
	}
	return fr.state.msgContainer.AllMessagesReceivedFor(fr.state.GetCurrentRound(), fr.instanceParams.operators)
}

func (fr *Instance) needToRunCurrentRound() bool {
	if !fr.instanceParams.isResharing() {
		return true // always run for new keygen
	}
	switch fr.state.GetCurrentRound() {
	case common.Preparation, common.Round2, common.KeygenOutput:
		return fr.instanceParams.inNewCommittee()
	case common.Round1:
		return fr.instanceParams.inOldCommittee()
	default:
		return false
	}
}

func (fr *Instance) validateSignedMessage(msg *dkg.SignedMessage) error {
	if msg.Message.Identifier != fr.instanceParams.identifier {
		return errors.New("got mismatching identifier")
	}

	found, operator, err := fr.config.GetStorage().GetDKGOperator(msg.Signer)
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

	addr := ethcommon.BytesToAddress(crypto.Keccak256(pk[1:])[12:])
	if addr != operator.ETHAddress {
		return errors.New("invalid signature")
	}
	return nil
}

func (fr *Instance) saveSignedMsg(msg *ProtocolMsg) (*dkg.SignedMessage, error) {
	bcastMessage, err := msg.ToSignedMessage(
		fr.instanceParams.identifier,
		fr.instanceParams.operatorID,
		fr.config.GetStorage(),
		fr.config.GetSigner(),
	)
	if err != nil {
		return nil, err
	}

	if _, err = fr.state.msgContainer.SaveMsg(fr.state.GetCurrentRound(), bcastMessage); err != nil {
		return nil, err
	}
	return bcastMessage, nil
}
