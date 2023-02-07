package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponCommit returns true if a quorum of commit messages was received.
// Assumes commit message is valid!
func (i *Instance) UponCommit(signedCommit *SignedMessage, commitMsgContainer *MsgContainer) (bool, []byte, *SignedMessage, error) {
	addMsg, err := commitMsgContainer.AddFirstMsgForSignerAndRound(signedCommit)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not add commit msg to container")
	}
	if !addMsg {
		return false, nil, nil, nil // UponCommit was already called
	}

	// calculate commit quorum and act upon it
	quorum, commitMsgs, err := commitQuorumForRoundRoot(i.State, commitMsgContainer, signedCommit.Message.Root, signedCommit.Message.Round)
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
func commitQuorumForRoundRoot(state *State, commitMsgContainer *MsgContainer, root [32]byte, round Round) (bool, []*SignedMessage, error) {
	signers, msgs := commitMsgContainer.LongestUniqueSignersForRoundAndRoot(round, root)
	return state.Share.HasQuorum(len(signers)), msgs, nil
}

func aggregateCommitMsgs(msgs []*SignedMessage, fullData []byte) (*SignedMessage, error) {
	if len(msgs) == 0 {
		return nil, errors.New("can't aggregate zero commit msgs")
	}

	var ret *SignedMessage
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

// didSendCommitForHeightAndRound returns true if sent commit msg for specific Height and round
/**
!exists m :: && m in current.messagesReceived
                            && m.Commit?
                            && var uPayload := m.commitPayload.unsignedPayload;
                            && uPayload.Height == |current.blockchain|
                            && uPayload.round == current.round
                            && recoverSignedCommitAuthor(m.commitPayload) == current.id
*/
func didSendCommitForHeightAndRound(state *State, commitMsgContainer *MsgContainer) bool {
	for _, msg := range commitMsgContainer.MessagesForRound(state.Round) {
		if msg.MatchedSigners([]types.OperatorID{state.Share.OperatorID}) {
			return true
		}
	}
	return false
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
func CreateCommit(state *State, config IConfig, root [32]byte) (*SignedMessage, error) {
	msg := &Message{
		MsgType:    CommitMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,

		Root: root,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing commit msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   *msg,
	}
	return signedMsg, nil
}

func baseCommitValidation(
	config IConfig,
	signedCommit *SignedMessage,
	height Height,
	operators []*types.Operator,
) error {
	if signedCommit.Message.MsgType != CommitMsgType {
		return errors.New("commit msg type is wrong")
	}
	if signedCommit.Message.Height != height {
		return errors.New("wrong msg height")
	}

	if err := signedCommit.Validate(); err != nil {
		return errors.Wrap(err, "signed commit invalid")
	}

	// verify signature
	if err := signedCommit.Signature.VerifyByOperators(signedCommit, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	return nil
}

func validateCommit(
	config IConfig,
	signedCommit *SignedMessage,
	height Height,
	round Round,
	proposedMsg *SignedMessage,
	operators []*types.Operator,
) error {
	if err := baseCommitValidation(config, signedCommit, height, operators); err != nil {
		return err
	}

	if len(signedCommit.Signers) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if signedCommit.Message.Round != round {
		return errors.New("wrong msg round")
	}

	if !bytes.Equal(proposedMsg.Message.Root[:], signedCommit.Message.Root[:]) {
		return errors.New("proposed data mistmatch")
	}

	return nil
}
