package qbft

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// uponPrepare process prepare message
// Assumes prepare message is valid!
func (i *Instance) uponPrepare(signedPrepare *types.SignedSSVMessage, prepareMsgContainer *MsgContainer) error {
	hasQuorumBefore := HasQuorum(i.State.Share, prepareMsgContainer.MessagesForRound(i.State.Round))

	addedMsg, err := prepareMsgContainer.AddFirstMsgForSignerAndRound(signedPrepare)
	if err != nil {
		return errors.Wrap(err, "could not add prepare msg to container")
	}
	if !addedMsg {
		return nil // uponPrepare was already called
	}

	if hasQuorumBefore {
		return nil // already moved to commit stage
	}

	if !HasQuorum(i.State.Share, prepareMsgContainer.MessagesForRound(i.State.Round)) {
		return nil // no quorum yet
	}

	proposalAcceptedMessage := &Message{}
	if err := proposalAcceptedMessage.Decode(i.State.ProposalAcceptedForCurrentRound.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "Could not decode Message from ProposalAcceptedForCurrentRound")
	}

	proposedRoot := proposalAcceptedMessage.Root

	i.State.LastPreparedValue = proposalAcceptedMessage.FullData
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
func getRoundChangeJustification(state *State, config IConfig, prepareMsgContainer *MsgContainer) ([]*types.SignedSSVMessage, error) {
	if state.LastPreparedValue == nil {
		return nil, nil
	}

	r, err := HashDataRoot(state.LastPreparedValue)
	if err != nil {
		return nil, errors.Wrap(err, "could not hash input data")
	}

	prepareMsgs := prepareMsgContainer.MessagesForRound(state.LastPreparedRound)
	ret := make([]*types.SignedSSVMessage, 0)
	for _, msg := range prepareMsgs {
		if err := validSignedPrepareForHeightRoundAndRoot(
			config,
			msg,
			state.Height,
			state.LastPreparedRound,
			r,
			state.Share.Committee,
			true,
		); err == nil {
			ret = append(ret, msg)
		}
	}

	if !HasQuorum(state.Share, ret) {
		return nil, nil
	}
	return ret, nil
}

// validSignedPrepareForHeightRoundAndRoot known in dafny spec as validSignedPrepareForHeightRoundAndDigest
// https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification
func validSignedPrepareForHeightRoundAndRoot(
	config IConfig,
	signedPrepare *types.SignedSSVMessage,
	height Height,
	round Round,
	root [32]byte,
	operators []*types.Operator,
	verifySignature bool) error {

	// Decode
	message := &Message{}
	if err := message.Decode(signedPrepare.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "Could not decode Prepare Message")
	}

	if message.MsgType != PrepareMsgType {
		return errors.New("prepare msg type is wrong")
	}
	if message.Height != height {
		return errors.New("wrong msg height")
	}
	if message.Round != round {
		return errors.New("wrong msg round")
	}

	if err := signedPrepare.Validate(); err != nil {
		return errors.Wrap(err, "prepareData invalid")
	}

	if !bytes.Equal(message.Root[:], root[:]) {
		return errors.New("proposed data mistmatch")
	}

	if len(signedPrepare.GetOperatorIDs()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if verifySignature {
		if err := types.VerifySignedSSVMessage(signedPrepare, operators); err != nil {
			return errors.Wrap(err, "msg signature invalid")
		}
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
func CreatePrepare(state *State, config IConfig, newRound Round, root [32]byte) (*Message, error) {
	return &Message{
		MsgType:    PrepareMsgType,
		Height:     state.Height,
		Round:      newRound,
		Identifier: state.ID,

		Root: root,
	}, nil
}
