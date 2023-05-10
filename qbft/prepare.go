package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// uponPrepare process prepare message
// Assumes prepare message is valid!
func (i *Instance) uponPrepare(
	signedPrepare *SignedMessage,
	prepareMsgContainer,
	commitMsgContainer *MsgContainer) error {

	addedMsg, err := prepareMsgContainer.AddFirstMsgForSignerAndRound(signedPrepare)
	if err != nil {
		return errors.Wrap(err, "could not add prepare msg to container")
	}
	if !addedMsg {
		return nil // uponPrepare was already called
	}

	if !HasQuorum(i.State.Share, prepareMsgContainer.MessagesForRound(i.State.Round)) {
		return nil // no quorum yet
	}

	if didSendCommitForHeightAndRound(i.State, commitMsgContainer) {
		return nil // already moved to commit stage
	}

	proposedRoot := i.State.ProposalAcceptedForCurrentRound.Message.Root

	i.State.LastPreparedValue = i.State.ProposalAcceptedForCurrentRound.FullData
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
func getRoundChangeJustification(state *State, config IConfig, prepareMsgContainer *MsgContainer) ([]*SignedMessage, error) {
	if state.LastPreparedValue == nil {
		return nil, nil
	}

	r, err := HashDataRoot(state.LastPreparedValue)
	if err != nil {
		return nil, errors.Wrap(err, "could not hash input data")
	}

	prepareMsgs := prepareMsgContainer.MessagesForRound(state.LastPreparedRound)
	ret := make([]*SignedMessage, 0)
	for _, msg := range prepareMsgs {
		if err := validSignedPrepareForHeightRoundAndRoot(
			config,
			msg,
			state.Height,
			state.LastPreparedRound,
			r,
			state.Share.Committee,
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
	signedPrepare *SignedMessage,
	height Height,
	round Round,
	root [32]byte,
	operators []*types.Operator) error {
	if signedPrepare.Message.MsgType != PrepareMsgType {
		return errors.New("prepare msg type is wrong")
	}
	if signedPrepare.Message.Height != height {
		return errors.New("wrong msg height")
	}
	if signedPrepare.Message.Round != round {
		return errors.New("wrong msg round")
	}

	if err := signedPrepare.Validate(); err != nil {
		return errors.Wrap(err, "prepareData invalid")
	}

	if !bytes.Equal(signedPrepare.Message.Root[:], root[:]) {
		return errors.New("proposed data mistmatch")
	}

	if len(signedPrepare.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if err := signedPrepare.Signature.VerifyByOperators(signedPrepare, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
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
func CreatePrepare(state *State, config IConfig, newRound Round, root [32]byte) (*SignedMessage, error) {
	msg := &Message{
		MsgType:    PrepareMsgType,
		Height:     state.Height,
		Round:      newRound,
		Identifier: state.ID,

		Root: root,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing prepare msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   *msg,
	}
	return signedMsg, nil
}
