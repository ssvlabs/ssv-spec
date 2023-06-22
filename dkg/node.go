package dkg

import (
	"encoding/hex"
	"fmt"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// Runners is a map of dkg runners mapped by dkg ID.
type Runners map[string]Runner

func (runners Runners) AddRunner(id RequestID, runner Runner) {
	runners[hex.EncodeToString(id[:])] = runner
}

// RunnerForID returns a Runner from the provided msg ID, or nil if not found
func (runners Runners) RunnerForID(id RequestID) Runner {
	return runners[hex.EncodeToString(id[:])]
}

func (runners Runners) Exists(id RequestID) bool {
	_, ok := runners[hex.EncodeToString(id[:])]
	return ok
}

func (runners Runners) DeleteRunner(id RequestID) {
	delete(runners, hex.EncodeToString(id[:]))
}

// Node is responsible for receiving and managing DKG session and messages
type Node struct {
	operator *Operator
	runners  Runners // runners holds all active running DKG runners
	config   *Config
}

func NewNode(operator *Operator, config *Config) *Node {
	return &Node{
		operator: operator,
		config:   config,
		runners:  make(Runners, 0),
	}
}

func (n *Node) newRunner(id RequestID, initMsg *Init) (Runner, error) {
	r := &runner{
		Operator:              n.operator,
		InitMsg:               initMsg,
		Identifier:            id,
		KeygenOutcome:         nil,
		DepositDataRoot:       nil,
		DepositDataSignatures: map[types.OperatorID]*PartialDepositData{},
		OutputMsgs:            map[types.OperatorID]*SignedOutput{},
		protocol:              n.config.KeygenProtocol(id, n.operator.OperatorID, n.config, initMsg),
		config:                n.config,
	}

	if err := r.protocol.Start(); err != nil {
		return nil, errors.Wrap(err, "could not start dkg protocol")
	}

	return r, nil
}

func (n *Node) newResharingRunner(id RequestID, reshareMsg *Reshare) (Runner, error) {
	reshareParams := &ReshareParams{}
	if inOldCommittee(n.operator.OperatorID, reshareMsg.OldOperatorIDs) {
		if err := reshareParams.LoadFromStorage(reshareMsg.ValidatorPK, n.config.Storage); err != nil {
			return nil, err
		}
	}

	r := &runner{
		Operator:              n.operator,
		ReshareMsg:            reshareMsg,
		Identifier:            id,
		KeygenOutcome:         nil,
		DepositDataRoot:       nil,
		DepositDataSignatures: map[types.OperatorID]*PartialDepositData{},
		OutputMsgs:            map[types.OperatorID]*SignedOutput{},
		protocol:              n.config.ReshareProtocol(id, n.operator.OperatorID, n.config, reshareMsg, reshareParams),
		config:                n.config,
	}

	if err := r.protocol.Start(); err != nil {
		return nil, errors.Wrap(err, "could not start resharing protocol")
	}

	return r, nil
}

func (n *Node) newSignatureRunner(id RequestID, keySign *KeySign) (Runner, error) {
	keygenOutput, err := n.config.GetStorage().GetKeyGenOutput(keySign.ValidatorPK)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret keyshare for given validator pk: %s", err.Error())
	}

	operators := make([]uint32, 0)
	for operatorID := range keygenOutput.OperatorPubKeys {
		operators = append(operators, uint32(operatorID))
	}

	keySign.Operators = operators
	keySign.Threshold = keygenOutput.Threshold
	keySign.SecretShare = keygenOutput.Share
	keySign.OperatorPublicKeyshares = keygenOutput.OperatorPubKeys

	r := &runner{
		Operator:              n.operator,
		KeySign:               keySign,
		Identifier:            id,
		KeygenOutcome:         nil,
		DepositDataRoot:       nil,
		DepositDataSignatures: map[types.OperatorID]*PartialDepositData{},
		OutputMsgs:            map[types.OperatorID]*SignedOutput{},
		protocol:              n.config.KeySign(id, n.operator.OperatorID, n.config, keySign),
		config:                n.config,
	}

	if err := r.protocol.Start(); err != nil {
		return nil, errors.Wrap(err, "could not start keysign protocol")
	}

	return r, nil
}

// ProcessMessage processes network Messages of all types
func (n *Node) ProcessMessage(msg *types.SSVMessage) error {
	if msg.MsgType != types.DKGMsgType {
		return errors.New("not a DKGMsgType")
	}
	signedMsg := &SignedMessage{}
	if err := signedMsg.Decode(msg.GetData()); err != nil {
		return errors.Wrap(err, "could not get dkg Message from network Messages")
	}

	if err := n.validateSignedMessage(signedMsg); err != nil {
		return errors.Wrap(err, "signed message doesn't pass validation")
	}

	switch signedMsg.Message.MsgType {
	case InitMsgType:
		return n.startNewDKGMsg(signedMsg)
	case ReshareMsgType:
		return n.startResharing(signedMsg)
	case KeySignMsgType:
		return n.startKeySign(signedMsg)
	case ProtocolMsgType:
		return n.processDKGMsg(signedMsg)
	case DepositDataMsgType:
		return n.processDKGMsg(signedMsg)
	case OutputMsgType:
		return n.processDKGMsg(signedMsg)
	default:
		return errors.New("unknown msg type")
	}
}

func (n *Node) validateSignedMessage(message *SignedMessage) error {
	if err := message.Validate(); err != nil {
		return errors.Wrap(err, "message invalid")
	}

	return nil
}

func (n *Node) startNewDKGMsg(message *SignedMessage) error {
	initMsg, err := n.validateInitMsg(message)
	if err != nil {
		return errors.Wrap(err, "could not start new dkg")
	}

	runner, err := n.newRunner(message.Message.Identifier, initMsg)
	if err != nil {
		return errors.Wrap(err, "could not start new dkg")
	}

	// add runner to runners
	n.runners.AddRunner(message.Message.Identifier, runner)

	return nil
}

func (n *Node) startResharing(message *SignedMessage) error {
	reshareMsg, err := n.validateReshareMsg(message)
	if err != nil {
		return errors.Wrap(err, "could not start resharing")
	}

	r, err := n.newResharingRunner(message.Message.Identifier, reshareMsg)
	if err != nil {
		return errors.Wrap(err, "could not start resharing")
	}

	// add runner to runners
	n.runners.AddRunner(message.Message.Identifier, r)

	return nil
}

func (n *Node) startKeySign(message *SignedMessage) error {
	keySign, err := n.validateKeySignMsg(message)
	if err != nil {
		return errors.Wrap(err, "could not start resharing")
	}

	r, err := n.newSignatureRunner(message.Message.Identifier, keySign)
	if err != nil {
		return errors.Wrap(err, "could not start resharing")
	}

	// add runner to runners
	n.runners.AddRunner(message.Message.Identifier, r)

	return nil
}

func (n *Node) validateInitMsg(message *SignedMessage) (*Init, error) {
	// validate identifier.GetEthAddress is the signer for message
	if err := message.Signature.ECRecover(message, n.config.SignatureDomainType, types.DKGSignatureType, message.Message.Identifier.GetETHAddress()); err != nil {
		return nil, errors.Wrap(err, "signed message invalid")
	}

	initMsg := &Init{}
	if err := initMsg.Decode(message.Message.Data); err != nil {
		return nil, errors.Wrap(err, "could not get dkg init Message from signed Messages")
	}

	if err := initMsg.Validate(); err != nil {
		return nil, errors.Wrap(err, "init message invalid")
	}

	// check instance not running already
	if n.runners.RunnerForID(message.Message.Identifier) != nil {
		return nil, errors.New("dkg started already")
	}

	return initMsg, nil
}

func (n *Node) validateReshareMsg(message *SignedMessage) (*Reshare, error) {
	// validate identifier.GetEthAddress is the signer for message
	if err := message.Signature.ECRecover(message, n.config.SignatureDomainType, types.DKGSignatureType, message.Message.Identifier.GetETHAddress()); err != nil {
		return nil, errors.Wrap(err, "signed message invalid")
	}

	reshareMsg := &Reshare{}
	if err := reshareMsg.Decode(message.Message.Data); err != nil {
		return nil, errors.Wrap(err, "could not get reshare Message from signed Messages")
	}

	if err := reshareMsg.Validate(); err != nil {
		return nil, errors.Wrap(err, "reshare message invalid")
	}

	// check instance not running already
	if n.runners.RunnerForID(message.Message.Identifier) != nil {
		return nil, errors.New("dkg started already")
	}

	return reshareMsg, nil
}

func (n *Node) validateKeySignMsg(message *SignedMessage) (*KeySign, error) {
	// validate identifier.GetEthAddress is the signer for message
	if err := message.Signature.ECRecover(message, n.config.SignatureDomainType, types.DKGSignatureType, message.Message.Identifier.GetETHAddress()); err != nil {
		return nil, errors.Wrap(err, "signed message invalid")
	}

	keySign := &KeySign{}
	if err := keySign.Decode(message.Message.Data); err != nil {
		return nil, errors.Wrap(err, "could not get validator exit msg params from signed message")
	}

	if err := keySign.Validate(); err != nil {
		return nil, errors.Wrap(err, "reshare message invalid")
	}

	// check instance not running already
	if n.runners.RunnerForID(message.Message.Identifier) != nil {
		return nil, errors.New("signature protocol started already")
	}

	return keySign, nil
}

func (n *Node) processDKGMsg(message *SignedMessage) error {
	if !n.runners.Exists(message.Message.Identifier) {
		return errors.New("could not find dkg runner")
	}

	if err := n.validateDKGMsg(message); err != nil {
		return errors.Wrap(err, "dkg msg not valid")
	}

	r := n.runners.RunnerForID(message.Message.Identifier)
	finished, err := r.ProcessMsg(message)
	if err != nil {
		return errors.Wrap(err, "could not process dkg message")
	}
	if finished {
		n.runners.DeleteRunner(message.Message.Identifier)
	}

	return nil
}

func (n *Node) validateDKGMsg(message *SignedMessage) error {

	// find signing operator and verify sig
	found, signingOperator, err := n.config.Storage.GetDKGOperator(message.Signer)
	if err != nil {
		return errors.Wrap(err, "can't fetch operator")
	}
	if !found {
		return errors.New("can't find operator")
	}
	if err := message.Signature.ECRecover(message, n.config.SignatureDomainType, types.DKGSignatureType, signingOperator.ETHAddress); err != nil {
		return errors.Wrap(err, "signed message invalid")
	}

	return nil
}

func (n *Node) GetConfig() *Config {
	return n.config
}
