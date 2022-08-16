package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponCommit returns true if a quorum of commit messages was received.
func (i *Instance) UponCommit(signedCommit *SignedMessage, commitMsgContainer *MsgContainer) (bool, []byte, *SignedMessage, error) {
	if err := validateCommit(
		i.State,
		i.config,
		signedCommit,
		i.State.Height,
		i.State.Share.Committee,
	); err != nil {
		return false, nil, nil, errors.Wrap(err, "commit msg invalid")
	}

	addMsg, err := commitMsgContainer.AddIfDoesntExist(signedCommit)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not add commit msg to container")
	}
	if !addMsg {
		return false, nil, nil, nil // UponCommit was already called
	}

	// calculate commit quorum and act upon it
	quorum, commitMsgs, err := commitQuorumForRoundValue(i.State, commitMsgContainer, signedCommit.Message.Data, signedCommit.Message.Round)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not calculate commit quorum")
	}
	if quorum {
		msgCommitData, err := signedCommit.Message.GetCommitData()
		if err != nil {
			return false, nil, nil, errors.Wrap(err, "could not get msg commit data")
		}

		agg, err := aggregateCommitMsgs(commitMsgs)
		if err != nil {
			return false, nil, nil, errors.Wrap(err, "could not aggregate commit msgs")
		}
		return true, msgCommitData.Data, agg, nil
	}
	return false, nil, nil, nil
}

// returns true if there is a quorum for the current round for this provided value
func commitQuorumForRoundValue(state *State, commitMsgContainer *MsgContainer, value []byte, round Round) (bool, []*SignedMessage, error) {
	signers, msgs := commitMsgContainer.LongestUniqueSignersForRoundAndValue(round, value)
	return state.Share.HasQuorum(len(signers)), msgs, nil
}

func aggregateCommitMsgs(msgs []*SignedMessage) (*SignedMessage, error) {
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
func CreateCommit(state *State, config IConfig, value []byte) (*SignedMessage, error) {
	commitData := &CommitData{
		Data: value,
	}
	dataByts, err := commitData.Encode()

	msg := &Message{
		MsgType:    CommitMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing commit msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}

func validateCommit(
	state *State,
	config IConfig,
	signedCommit *SignedMessage,
	height Height,
	operators []*types.Operator,
) error {
	if signedCommit.Message.MsgType != CommitMsgType {
		return errors.New("commit msg type is wrong")
	}
	if signedCommit.Message.Height != height {
		return errors.New("commit Height is wrong")
	}

	msgCommitData, err := signedCommit.Message.GetCommitData()
	if err != nil {
		return errors.Wrap(err, "could not get msg commit data")
	}
	if err := msgCommitData.Validate(); err != nil {
		return errors.Wrap(err, "msgCommitData invalid")
	}

	if err := signedCommit.Signature.VerifyByOperators(signedCommit, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "commit msg signature invalid")
	}

	if isDecidedMsg(state, signedCommit) {
		if err := config.GetValueCheckF()(msgCommitData.Data); err != nil {
			return errors.Wrap(err, "invalid decided data")
		}
	}

	return nil
}

// returns true if signed commit has all quorum sigs
func isDecidedMsg(state *State, signedCommit *SignedMessage) bool {
	return state.Share.HasQuorum(len(signedCommit.Signers))
}
