package qbft

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// UponCommit returns true if a quorum of commit messages was received.
// Assumes commit message is valid!
func (i *Instance) UponCommit(msg *ProcessingMessage, commitMsgContainer *MsgContainer) (bool, []byte, *types.SignedSSVMessage, error) {

	addMsg, err := commitMsgContainer.AddFirstMsgForSignerAndRound(msg)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not add commit msg to container")
	}
	if !addMsg {
		return false, nil, nil, nil // UponCommit was already called
	}

	// calculate commit quorum and act upon it
	quorum, commitMsgs, err := commitQuorumForRoundRoot(i.State, commitMsgContainer, msg.QBFTMessage.Root, msg.QBFTMessage.Round)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not calculate commit quorum")
	}
	if quorum {
		fullData := i.State.ProposalAcceptedForCurrentRound.SignedMessage.FullData /* must have value there, checked on validateCommit */

		agg, err := aggregateCommitMsgs(commitMsgs, fullData)
		if err != nil {
			return false, nil, nil, errors.Wrap(err, "could not aggregate commit msgs")
		}
		return true, fullData, agg, nil
	}
	return false, nil, nil, nil
}

// returns true if there is a quorum for the current round for this provided value
func commitQuorumForRoundRoot(state *State, commitMsgContainer *MsgContainer, root [32]byte, round Round) (bool, []*ProcessingMessage, error) {
	signers, msgs := commitMsgContainer.LongestUniqueSignersForRoundAndRoot(round, root)
	return state.CommitteeMember.HasQuorum(len(signers)), msgs, nil
}

func aggregateCommitMsgs(msgs []*ProcessingMessage, fullData []byte) (*types.SignedSSVMessage, error) {
	if len(msgs) == 0 {
		return nil, errors.New("can't aggregate zero commit msgs")
	}

	var ret *types.SignedSSVMessage
	for _, m := range msgs {
		if ret == nil {
			ret = m.SignedMessage.DeepCopy()
		} else {
			if err := ret.Aggregate(m.SignedMessage); err != nil {
				return nil, errors.Wrap(err, "could not aggregate commit msg")
			}
		}
	}
	ret.FullData = fullData
	return ret, nil
}

// CreateCommit
/**
Commit(
                    signCommit(
                        UnsignedCommit(
                            |current.blockchain|,
                            current.round,
                            signHash(hashBlockForCommitSeal(proposedBlock), current.id),
                            digest(proposedBlock)),
                            current.id
                        )
                    );
*/
func CreateCommit(state *State, config IConfig, root [32]byte) (*types.SignedSSVMessage, error) {
	msg := &Message{
		MsgType:    CommitMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,

		Root: root,
	}
	return Sign(msg, state.CommitteeMember.OperatorID, config.GetOperatorSigner())
}

func baseCommitValidationIgnoreSignature(
	msg *ProcessingMessage,
	height Height,
	operators []*types.Operator,
) error {

	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "signed commit invalid")
	}

	if msg.QBFTMessage.MsgType != CommitMsgType {
		return errors.New("commit msg type is wrong")
	}
	if msg.QBFTMessage.Height != height {
		return errors.New("wrong msg height")
	}

	if !msg.SignedMessage.CheckSignersInCommittee(operators) {
		return errors.New("signer not in committee")
	}

	return nil
}

func baseCommitValidationVerifySignature(
	config IConfig,
	msg *ProcessingMessage,
	height Height,
	operators []*types.Operator) error {

	if err := baseCommitValidationIgnoreSignature(msg, height, operators); err != nil {
		return err
	}

	// verify signature
	if err := config.GetSignatureVerifier().Verify(msg.SignedMessage, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	return nil
}

func validateCommit(
	msg *ProcessingMessage,
	height Height,
	round Round,
	proposedMsg *ProcessingMessage,
	operators []*types.Operator,
) error {
	if err := baseCommitValidationIgnoreSignature(msg, height, operators); err != nil {
		return err
	}

	if len(msg.SignedMessage.OperatorIDs) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if msg.QBFTMessage.Round != round {
		return errors.New("wrong msg round")
	}

	if !bytes.Equal(proposedMsg.QBFTMessage.Root[:], msg.QBFTMessage.Root[:]) {
		return errors.New("proposed data mismatch")
	}

	return nil
}
