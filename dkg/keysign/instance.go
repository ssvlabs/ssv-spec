package keysign

import (
	"fmt"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/coinbase/kryptology/pkg/core/curves"
	ethcommon "github.com/ethereum/go-ethereum/common"
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

type InstanceParams struct {
	operatorID types.OperatorID
	identifier dkg.RequestID

	validatorPK types.ValidatorPK
	threshold   uint64
	operators   []uint32

	Key                     *bls.SecretKey
	OperatorPublicKeyshares map[types.OperatorID]*bls.PublicKey
	SigningRoot             []byte
}

type Instance struct {
	state          *State
	config         dkg.IConfig
	instanceParams InstanceParams
}

// NewSignature creates a protocol instance for computing a new signature
func NewSignature(
	requestID dkg.RequestID,
	operatorID types.OperatorID,
	config dkg.IConfig,
	init *dkg.KeySign,
) (dkg.Protocol, error) {

	keygenOutput, err := config.GetStorage().GetKeyGenOutput(init.ValidatorPK)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret keyshare for given validator pk: %s", err.Error())
	}

	operators := make([]uint32, 0)
	for operatorID, _ := range keygenOutput.OperatorPubKeys {
		operators = append(operators, uint32(operatorID))
	}

	instanceParams := InstanceParams{
		identifier: requestID,
		operatorID: operatorID,

		validatorPK: init.ValidatorPK,
		threshold:   keygenOutput.Threshold,
		operators:   operators,

		Key:                     keygenOutput.Share,
		OperatorPublicKeyshares: keygenOutput.OperatorPubKeys,
		SigningRoot:             init.SigningRoot,
	}
	return &Instance{
		config:         config,
		state:          initState(),
		instanceParams: instanceParams,
	}, nil
}

// Start ...
func (instance *Instance) Start() error {
	instance.state.roundTimer.OnTimeout(instance.UponRoundTimeout)
	instance.state.SetCurrentRound(common.Preparation)
	instance.state.roundTimer.StartRoundTimeoutTimer(common.ProtocolRound(instance.state.GetCurrentRound()))

	partialSignature := instance.instanceParams.Key.SignByte(instance.instanceParams.SigningRoot)

	msg := &ProtocolMsg{
		Round: common.Preparation,
		PreparationMessage: &PreparationMessage{
			PartialSignature: partialSignature.Serialize(),
		},
	}

	bcastMsg, err := instance.saveSignedMsg(msg)
	if err != nil {
		return err
	}
	return instance.config.GetNetwork().BroadcastDKGMessage(bcastMsg)
}

// ProcessMsg ...
func (instance *Instance) ProcessMsg(msg *dkg.SignedMessage) (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {
	// validated signed message
	if err := instance.validateSignedMessage(msg); err != nil {
		return false, nil, errors.Wrap(err, "failed to Validate signed message")
	}

	// decodes and validates protocol message
	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "failed to decode protocol msg")
	}
	if err := protocolMessage.Validate(); err != nil {
		return true, nil, fmt.Errorf("signature protocol failed: invalid protocol msg: %s", err)
	}

	_, err = instance.state.msgContainer.SaveMsg(protocolMessage.Round, msg)
	if err != nil {
		return true, nil, err // message already exists or err with saving msg
	}

	switch protocolMessage.Round {
	case common.Preparation:
		return instance.processRound1()
	case common.Round1:
		return instance.processKeysignOutput()
	default:
		return true, nil, dkg.ErrInvalidRound{}
	}
}

func (instance *Instance) canProceedThisRound() bool {
	switch instance.state.GetCurrentRound() {
	case common.Preparation:
		instance.state.msgContainer.AllMessagesReceivedUpto(instance.state.GetCurrentRound(), instance.instanceParams.operators, instance.instanceParams.threshold)
	case common.Round1:
		instance.state.msgContainer.AllMessagesReceivedFor(instance.state.GetCurrentRound(), instance.instanceParams.operators)
	}
	return false
}

func (instance *Instance) validateSignedMessage(msg *dkg.SignedMessage) error {
	if msg.Message.Identifier != instance.instanceParams.identifier {
		return errors.New("got mismatching identifier")
	}

	found, operator, err := instance.config.GetStorage().GetDKGOperator(msg.Signer)
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

func (instance *Instance) saveSignedMsg(msg *ProtocolMsg) (*dkg.SignedMessage, error) {
	bcastMessage, err := msg.ToSignedMessage(
		instance.instanceParams.identifier,
		instance.instanceParams.operatorID,
		instance.config.GetStorage(),
		instance.config.GetSigner(),
	)
	if err != nil {
		return nil, err
	}

	if _, err = instance.state.msgContainer.SaveMsg(common.ProtocolRound(instance.state.GetCurrentRound()), bcastMessage); err != nil {
		return nil, err
	}
	return bcastMessage, nil
}
