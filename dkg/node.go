package dkg

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// Runners is a map of dkg runners mapped by dkg ID.
type Runners map[string]*Runner

func (runners Runners) AddRunner(id RequestID, runner *Runner) {
	runners[hex.EncodeToString(id[:])] = runner
}

// RunnerForID returns a Runner from the provided msg ID, or nil if not found
func (runners Runners) RunnerForID(id RequestID) *Runner {
	return runners[hex.EncodeToString(id[:])]
}

func (runners Runners) DeleteRunner(id RequestID) {
	delete(runners, hex.EncodeToString(id[:]))
}

type Node struct {
	operator *Operator
	// runners holds all active running DKG runners
	runners Runners
	config  *Config
}

func NewNode(operator *Operator, config *Config) *Node {
	return &Node{
		operator: operator,
		config:   config,
		runners:  make(Runners, 0),
	}
}

func (n *Node) newRunner(id RequestID, initMsg *Init) (*Runner, error) {

	var i uint16
	for i0, id := range initMsg.OperatorIDs {
		if id == n.operator.OperatorID {
			i = uint16(i0) + 1
		}
	}
	if i == 0 {
		return nil, errors.New("invalid request")
	}

	runner := &Runner{
		Operator:          n.operator,
		InitMsg:           initMsg,
		Identifier:        id,
		I:                 i,
		PartialSignatures: map[types.OperatorID][]byte{},
		DepositDataSignatures: map[types.OperatorID]*PartialDepositData{},
		config:            n.config,
		keygenSubProtocol: n.config.Protocol(initMsg, n.operator.OperatorID, id),
	}

	return runner, nil
}

// ProcessMessage processes network Messages of all types
func (n *Node) ProcessMessage(msg *types.SSVMessage) error {
	signedMsg := &SignedMessage{}
	if err := signedMsg.Decode(msg.GetData()); err != nil {
		return errors.Wrap(err, "could not get dkg Message from network Messages")
	}

	if err := n.validateSignedMessage(signedMsg); err != nil {
		return errors.Wrap(err, "signed message doesn't pass validation")
	}

	switch signedMsg.Message.Header.MsgType {
	case int32(InitMsgType):
		return n.startNewDKGMsg(signedMsg)
	case int32(ProtocolMsgType):
		return n.processDKGMsg(signedMsg)
	case int32(DepositDataMsgType):
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

	runner, err := n.newRunner(message.Message.Header.RequestID(), initMsg)
	if err != nil {
		return errors.Wrap(err, "could not start new dkg")
	}

	// add runner to runners
	n.runners.AddRunner(message.Message.Header.RequestID(), runner)
	err = runner.Start()
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) validateInitMsg(message *SignedMessage) (*Init, error) {
	// validate identifier.GetEthAddress is the signer for message
	if err := message.Signature.ECRecover(message, n.config.SignatureDomainType, types.DKGSignatureType, RequestID(message.Message.Header.RequestID()).GetETHAddress()); err != nil {
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
	if n.runners.RunnerForID(message.Message.Header.RequestID()) != nil {
		return nil, errors.New("dkg started already")
	}

	return initMsg, nil
}

func (n *Node) processDKGMsg(message *SignedMessage) error {
	runner, err := n.validateDKGMsg(message)
	if err != nil {
		return errors.Wrap(err, "dkg msg not valid")
	}

	finished, output, err := runner.ProcessMsg(message)
	if err != nil {
		return errors.Wrap(err, "could not process dkg message")
	}

	if finished {
		if err := n.config.Network.StreamDKGOutput(output); err != nil {
			return errors.Wrap(err, "failed to stream dkg output")
		}
		n.runners.DeleteRunner(message.Message.Header.RequestID())
	}

	return nil
}

func (n *Node) validateDKGMsg(message *SignedMessage) (*Runner, error) {
	runner := n.runners.RunnerForID(message.Message.Header.RequestID())
	if runner == nil {
		return nil, errors.New("could not find dkg runner")
	}

	// find signing operator and verify sig
	found, signingOperator, err := n.config.Storage.GetDKGOperator(message.Signer)
	if err != nil {
		return nil, errors.Wrap(err, "can't fetch operator")
	}
	if !found {
		return nil, errors.New("can't find operator")
	}
	if err := message.Signature.ECRecover(message, n.config.SignatureDomainType, types.DKGSignatureType, signingOperator.ETHAddress); err != nil {
		return nil, errors.Wrap(err, "signed message invalid")
	}

	return runner, nil
}
