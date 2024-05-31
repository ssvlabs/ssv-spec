package qbft

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// uponPrepare process prepare message
// Assumes prepare message is valid!
func (i *Instance) uponPrepare(signedPrepare *types.SignedSSVMessage, prepareMsgContainer *MsgContainer) error {
	hasQuorumBefore := HasQuorum(i.State.SharedValidator, prepareMsgContainer.MessagesForRound(i.State.Round))

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

	if !HasQuorum(i.State.SharedValidator, prepareMsgContainer.MessagesForRound(i.State.Round)) {
		return nil // no quorum yet
	}

	proposalMsgAccepted, err := DecodeMessage(i.State.ProposalAcceptedForCurrentRound.SSVMessage.Data)
	if err != nil {
		return err
	}

	proposedRoot := proposalMsgAccepted.Root

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
		if err := validSignedPrepareForHeightRoundAndRootIgnoreSignature(
			msg,
			state.Height,
			state.LastPreparedRound,
			r,
			state.SharedValidator.Committee,
		); err == nil {
			ret = append(ret, msg)
		}
	}

	if !HasQuorum(state.SharedValidator, ret) {
		return nil, nil
	}
	return ret, nil
}

// validSignedPrepareForHeightRoundAndRoot known in dafny spec as validSignedPrepareForHeightRoundAndDigest
// https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification
func validSignedPrepareForHeightRoundAndRootIgnoreSignature(
	signedPrepare *types.SignedSSVMessage,
	height Height,
	round Round,
	root [32]byte,
	operators []*types.ValidatorShare) error {

	msg, err := DecodeMessage(signedPrepare.SSVMessage.Data)
	if err != nil {
		return err
	}

	if msg.MsgType != PrepareMsgType {
		return errors.New("prepare msg type is wrong")
	}
	if msg.Height != height {
		return errors.New("wrong msg height")
	}
	if msg.Round != round {
		return errors.New("wrong msg round")
	}

	if err := signedPrepare.Validate(); err != nil {
		return errors.Wrap(err, "prepareData invalid")
	}

	if !bytes.Equal(msg.Root[:], root[:]) {
		return errors.New("proposed data mismatch")
	}

	if len(signedPrepare.GetOperatorIDs()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if !signedPrepare.CheckSignersInCommittee(operators) {
		return errors.New("signer not in committee")
	}

	return nil
}

func validSignedPrepareForHeightRoundAndRootVerifySignature(
	config IConfig,
	signedPrepare *types.SignedSSVMessage,
	height Height,
	round Round,
	root [32]byte,
	operators []*types.ValidatorShare) error {

	if err := validSignedPrepareForHeightRoundAndRootIgnoreSignature(signedPrepare, height, round, root, operators); err != nil {
		return err
	}

	// Verify signature
	if err := config.GetSignatureVerifier().Verify(signedPrepare, operators); err != nil {
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

	return MessageToSignedSSVMessage(msg, state.SharedValidator.OwnValidatorShare.OperatorID, config.GetOperatorSigner())
}
