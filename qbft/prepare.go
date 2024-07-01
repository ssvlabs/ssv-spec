package qbft

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// uponPrepare process prepare message
// Assumes prepare message is valid!
func (i *Instance) uponPrepare(msg *ProcessingMessage, prepareMsgContainer *MsgContainer) error {
	hasQuorumBefore := HasQuorum(i.State.CommitteeMember, prepareMsgContainer.MessagesForRound(i.State.Round))

	addedMsg, err := prepareMsgContainer.AddFirstMsgForSignerAndRound(msg)
	if err != nil {
		return errors.Wrap(err, "could not add prepare msg to container")
	}
	if !addedMsg {
		return nil // uponPrepare was already called
	}

	if hasQuorumBefore {
		return nil // already moved to commit stage
	}

	if !HasQuorum(i.State.CommitteeMember, prepareMsgContainer.MessagesForRound(i.State.Round)) {
		return nil // no quorum yet
	}

	proposedRoot := i.State.ProposalAcceptedForCurrentRound.QBFTMessage.Root

	i.State.LastPreparedValue = i.State.ProposalAcceptedForCurrentRound.SignedMessage.FullData
	i.State.LastPreparedRound = i.State.Round

	commitMsg, err := CreateCommit(i.State, i.config, proposedRoot)
	if err != nil {
		return errors.Wrap(err, "could not create commit msg")
	}

	if err := i.Broadcast(commitMsg); err != nil {
		return errors.Wrap(err, "failed to broadcast commit message")
	}

	return nil
}

// getRoundChangeJustification returns the round change justification for the current round.
// The justification is a quorum of signed prepare messages that agree on state.LastPreparedValue
func getRoundChangeJustification(state *State, config IConfig, prepareMsgContainer *MsgContainer) ([]*ProcessingMessage, error) {
	if state.LastPreparedValue == nil {
		return nil, nil
	}

	r, err := HashDataRoot(state.LastPreparedValue)
	if err != nil {
		return nil, errors.Wrap(err, "could not hash input data")
	}

	prepareMsgs := prepareMsgContainer.MessagesForRound(state.LastPreparedRound)
	ret := make([]*ProcessingMessage, 0)
	for _, msg := range prepareMsgs {
		if err := validSignedPrepareForHeightRoundAndRootIgnoreSignature(
			msg,
			state.Height,
			state.LastPreparedRound,
			r,
			state.CommitteeMember.Committee,
		); err == nil {
			ret = append(ret, msg)
		}
	}

	if !HasQuorum(state.CommitteeMember, ret) {
		return nil, nil
	}
	return ret, nil
}

// validSignedPrepareForHeightRoundAndRoot known in dafny spec as validSignedPrepareForHeightRoundAndDigest
// https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification
func validSignedPrepareForHeightRoundAndRootIgnoreSignature(
	msg *ProcessingMessage,
	height Height,
	round Round,
	root [32]byte,
	operators []*types.Operator) error {

	if msg.QBFTMessage.MsgType != PrepareMsgType {
		return errors.New("prepare msg type is wrong")
	}
	if msg.QBFTMessage.Height != height {
		return errors.New("wrong msg height")
	}
	if msg.QBFTMessage.Round != round {
		return errors.New("wrong msg round")
	}

	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "prepareData invalid")
	}

	if !bytes.Equal(msg.QBFTMessage.Root[:], root[:]) {
		return errors.New("proposed data mismatch")
	}

	if len(msg.SignedMessage.OperatorIDs) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if !msg.SignedMessage.CheckSignersInCommittee(operators) {
		return errors.New("signer not in committee")
	}

	return nil
}

func validSignedPrepareForHeightRoundAndRootVerifySignature(
	config IConfig,
	msg *ProcessingMessage,
	height Height,
	round Round,
	root [32]byte,
	operators []*types.Operator) error {

	if err := validSignedPrepareForHeightRoundAndRootIgnoreSignature(msg, height, round, root, operators); err != nil {
		return err
	}

	// Verify signature
	if err := config.GetSignatureVerifier().Verify(msg.SignedMessage, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	return nil
}

// CreatePrepare
/**
Prepare(
                    signPrepare(
                        UnsignedPrepare(
                            |current.blockchain|,
                            newRound,
                            digest(m.proposedBlock)),
                        current.id
                        )
                );
*/
func CreatePrepare(state *State, config IConfig, newRound Round, root [32]byte) (*types.SignedSSVMessage, error) {
	msg := &Message{
		MsgType:    PrepareMsgType,
		Height:     state.Height,
		Round:      newRound,
		Identifier: state.ID,

		Root: root,
	}

	return Sign(msg, state.CommitteeMember.OperatorID, config.GetOperatorSigner())
}
