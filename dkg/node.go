package dkg

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/dkg/stubdkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// Runners is a map of dkg runners mapped by dkg ID.
type Runners map[string]*Runner

func (runners Runners) AddRunner(id types.MessageID, runner *Runner) {
	runners[hex.EncodeToString(id)] = runner
}

// RunnerForID returns a Runner from the provided msg ID, or nil if not found
func (runners Runners) RunnerForID(msgID types.MessageID) *Runner {
	panic("implement")
}

type Node struct {
	OperatorID types.OperatorID
	// KeyPK signing public key
	KeyPK *ecdsa.PublicKey

	runners Runners
	network Network
	signer  types.DKGSigner
}

// ProcessMessage processes network Messages of all types
func (n *Node) ProcessMessage(msg *types.SSVMessage) error {
	// TODO validate msg

	signedMsg := &SignedMessage{}
	if err := signedMsg.Decode(msg.GetData()); err != nil {
		return errors.Wrap(err, "could not get dkg Message from network Messages")
	}

	switch signedMsg.Message.MsgType {
	case InitMsgType:
		return n.startNewDKGMsg(signedMsg)
	case ProtocolMsgType:
		return n.processDKGMsg(signedMsg)
	default:
		return errors.New("unknown msg type")
	}
}

func (n *Node) startNewDKGMsg(message *SignedMessage) error {
	initMsg := &Init{}
	if err := initMsg.Decode(message.Message.Data); err != nil {
		return errors.Wrap(err, "could not get dkg init Message from signed Messages")
	}

	// TODO - validate message
	// check instance not running already
	if n.runners.RunnerForID(message.Message.Identifier) != nil {
		return errors.New("dkg started already")
	}

	runner, err := NewRunner(initMsg, &Config{
		Protocol:   stubdkg.New(n.network, n.OperatorID, message.Message.Identifier),
		Network:    n.network,
		OperatorID: n.OperatorID,
		Identifier: message.Message.Identifier,
		Signer:     n.signer,
	})
	if err != nil {
		return errors.Wrap(err, "could not start new dkg")
	}

	// add runner to runners
	n.runners.AddRunner(message.Message.Identifier, runner)

	return nil
}

func (n *Node) processDKGMsg(message *SignedMessage) error {
	runner := n.runners.RunnerForID(message.Message.Identifier)
	if runner == nil {
		return errors.New("could not find dkg runner")
	}

	finished, output, err := runner.ProcessMsg(message)
	if err != nil {
		return errors.Wrap(err, "could not process dkg message")
	}

	if finished {
		if err := n.network.StreamDKGOutput(output); err != nil {
			return errors.Wrap(err, "failed to stream dkg output")
		}
	}

	return nil
}
