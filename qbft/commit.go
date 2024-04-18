package qbft

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponCommit returns true if a quorum of commit messages was received.
// Assumes commit message is valid!
func (i *Instance) UponCommit(signedCommit *types.SignedSSVMessage, commitMsgContainer *MsgContainer) (bool, []byte, *types.SignedSSVMessage, error) {
	// Decode qbft message
	msg, err := DecodeMessage(signedCommit.SSVMessage.Data)
	if err != nil {
		return false, nil, nil, err
	}

	addMsg, err := commitMsgContainer.AddFirstMsgForSignerAndRound(signedCommit)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not add commit msg to container")
	}
	if !addMsg {
		return false, nil, nil, nil // UponCommit was already called
	}

	// calculate commit quorum and act upon it
	quorum, commitMsgs, err := commitQuorumForRoundRoot(i.State, commitMsgContainer, msg.Root, msg.Round)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not calculate commit quorum")
	}
	if quorum {
		fullData := i.State.ProposalAcceptedForCurrentRound.FullData /* must have value there, checked on validateCommit */

		agg, err := aggregateCommitMsgs(commitMsgs, fullData)
		if err != nil {
			return false, nil, nil, errors.Wrap(err, "could not aggregate commit msgs")
		}
		return true, fullData, agg, nil
	}
	return false, nil, nil, nil
}

// returns true if there is a quorum for the current round for this provided value
func commitQuorumForRoundRoot(state *State, commitMsgContainer *MsgContainer, root [32]byte, round Round) (bool, []*types.SignedSSVMessage, error) {
	signers, msgs := commitMsgContainer.LongestUniqueSignersForRoundAndRoot(round, root)
	return state.Share.HasQuorum(len(signers)), msgs, nil
}

func aggregateCommitMsgs(msgs []*types.SignedSSVMessage, fullData []byte) (*types.SignedSSVMessage, error) {
	if len(msgs) == 0 {
		return nil, errors.New("can't aggregate zero commit msgs")
	}

	var ret *types.SignedSSVMessage
	for _, m := range msgs {
		if ret == nil {
			ret = m.DeepCopy()
		} else {
			if err := ret.Aggregate(m); err != nil {
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
	return MessageToSignedSSVMessage(msg, state.Share.OperatorID, config.GetOperatorSigner())
}

func baseCommitValidationIgnoreSignature(
	signedCommit *types.SignedSSVMessage,
	height Height,
	operators []*types.Operator,
) error {

	if err := signedCommit.Validate(); err != nil {
		return errors.Wrap(err, "signed commit invalid")
	}

	msg, err := DecodeMessage(signedCommit.SSVMessage.Data)
	if err != nil {
		return err
	}

	if msg.MsgType != CommitMsgType {
		return errors.New("commit msg type is wrong")
	}
	if msg.Height != height {
		return errors.New("wrong msg height")
	}

	if !signedCommit.CheckSignersInCommittee(operators) {
		return errors.New("signer not in committee")
	}

	return nil
}

func baseCommitValidationVerifySignature(
	config IConfig,
	signedCommit *types.SignedSSVMessage,
	height Height,
	operators []*types.Operator,
) error {

	if err := baseCommitValidationIgnoreSignature(signedCommit, height, operators); err != nil {
		return err
	}

	// verify signature
	if err := config.GetSignatureVerifier().Verify(signedCommit, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	return nil
}

func validateCommit(
	signedCommit *types.SignedSSVMessage,
	height Height,
	round Round,
	proposedSignedMsg *types.SignedSSVMessage,
	operators []*types.Operator,
) error {
	if err := baseCommitValidationIgnoreSignature(signedCommit, height, operators); err != nil {
		return err
	}

	if len(signedCommit.GetOperatorIDs()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	msg, err := DecodeMessage(signedCommit.SSVMessage.Data)
	if err != nil {
		return err
	}

	if msg.Round != round {
		return errors.New("wrong msg round")
	}

	proposedMsg, err := DecodeMessage(proposedSignedMsg.SSVMessage.Data)
	if err != nil {
		return err
	}

	if !bytes.Equal(proposedMsg.Root[:], msg.Root[:]) {
		return errors.New("proposed data mistmatch")
	}

	return nil
}
