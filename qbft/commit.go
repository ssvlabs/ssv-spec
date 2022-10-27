package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponCommit returns true if a quorum of commit messages was received.
func (i *Instance) UponCommit(
	signedCommit *SignedMessage,
	commitMsgContainer *MsgContainer,
) (bool, []byte, *SignedMessage, error) {
	if i.State.ProposalAcceptedForCurrentRound == nil {
		return false, nil, nil, errors.New("did not receive proposal for this round")
	}

	if err := validateCommit(
		i.config,
		signedCommit,
		i.State.Height,
		i.State.Round,
		i.State.ProposalAcceptedForCurrentRound.Message.Input.Root,
		i.State.Share.Committee,
	); err != nil {
		return false, nil, nil, errors.Wrap(err, "commit msg invalid")
	}

	addMsg, err := commitMsgContainer.AddFirstMsgForSignerAndRound(signedCommit)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not add commit msg to container")
	}
	if !addMsg {
		return false, nil, nil, nil // UponCommit was already called
	}

	// calculate commit quorum and act upon it
	quorum, commitMsgs, err := commitQuorumForRoundValue(i.State, commitMsgContainer, signedCommit.Message.Input.Root[:], signedCommit.Message.Round)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not calculate commit quorum")
	}
	if quorum {
		agg, err := aggregateCommitMsgs(commitMsgs, i.State.ProposalAcceptedForCurrentRound.Message.Input)
		if err != nil {
			return false, nil, nil, errors.Wrap(err, "could not aggregate commit msgs")
		}
		return true, i.State.ProposalAcceptedForCurrentRound.Message.Input.Source, agg, nil
	}
	return false, nil, nil, nil
}

// returns true if there is a quorum for the current round for this provided value
func commitQuorumForRoundValue(state *State, commitMsgContainer *MsgContainer, value []byte, round Round) (bool, []*SignedMessage, error) {
	signers, msgs := commitMsgContainer.LongestUniqueSignersForRoundAndValue(round, value)
	return state.Share.HasQuorum(len(signers)), msgs, nil
}

func aggregateCommitMsgs(msgs []*SignedMessage, acceptedProposalData *Data) (*SignedMessage, error) {
	if len(msgs) == 0 {
		return nil, errors.New("can't aggregate zero commit msgs")
	}

	var ret *SignedMessage
	for _, m := range msgs {
		if ret == nil {
			ret = m.DeepCopy(acceptedProposalData)
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
func CreateCommit(state *State, config IConfig, value [32]byte) (*SignedMessage, error) {
	msg := &Message{
		Height: state.Height,
		Round:  state.Round,
		Input: &Data{
			Root:   value,
			Source: nil,
		},
	}

	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing commit msg")
	}

	return &SignedMessage{
		Message:   msg,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Signature: sig,
	}, nil
}

func baseCommitValidation(
	config IConfig,
	signedCommit *SignedMessage,
	height Height,
	operators []*types.Operator,
) error {
	if signedCommit.Message.Height != height {
		return errors.New("commit Height is wrong")
	}

	if err := signedCommit.Signature.VerifyByOperators(signedCommit, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "commit msg signature invalid")
	}

	return nil
}

func validateCommit(
	config IConfig,
	signedCommit *SignedMessage,
	height Height,
	round Round,
	inputRoot [32]byte,
	operators []*types.Operator,
) error {
	if err := baseCommitValidation(config, signedCommit, height, operators); err != nil {
		return errors.Wrap(err, "invalid commit msg")
	}

	if len(signedCommit.Signers) != 1 {
		return errors.New("commit msgs allow 1 signer")
	}

	if signedCommit.Message.Round != round {
		return errors.New("commit round is wrong")
	}

	if !bytes.Equal(signedCommit.Message.Input.Root[:], inputRoot[:]) {
		return errors.New("proposed data different than commit msg data")
	}

	return nil
}
