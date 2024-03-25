package qbft

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponCommit returns true if a quorum of commit messages was received.
// Assumes commit message is valid!
func (i *Instance) UponCommit(signedCommit *types.SignedSSVMessage, commitMsgContainer *MsgContainer) (bool, []byte, *types.SignedSSVMessage, error) {
	addMsg, err := commitMsgContainer.AddFirstMsgForSignerAndRound(signedCommit)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not add commit msg to container")
	}
	if !addMsg {
		return false, nil, nil, nil // UponCommit was already called
	}

	// Decode
	message := &Message{}
	if err := message.Decode(signedCommit.SSVMessage.Data); err != nil {
		return false, nil, nil, errors.Wrap(err, "could not decode Commit Message")
	}

	// calculate commit quorum and act upon it
	quorum, commitMsgs, err := commitQuorumForRoundRoot(i.State, commitMsgContainer, message.Root, message.Round)
	if err != nil {
		return false, nil, nil, errors.Wrap(err, "could not calculate commit quorum")
	}
	if quorum {

		// Decode
		proposalMessage := &Message{}
		if err := proposalMessage.Decode(i.State.ProposalAcceptedForCurrentRound.SSVMessage.Data); err != nil {
			return false, nil, nil, errors.Wrap(err, "could not decode ProposalAcceptedForCurrentRound Message")
		}

		fullData := proposalMessage.FullData /* must have value there, checked on validateCommit */

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

	// Create signers and signatures list
	signers := make([]types.OperatorID, 0)
	signatures := make([][]byte, 0)

	for _, msg := range msgs {
		signers = append(signers, msg.OperatorID...)
		signatures = append(signatures, msg.Signature...)
	}

	// Insert FullData into Message
	sampleSSVMessage := msgs[0].SSVMessage
	message := &Message{}
	if err := message.Decode(sampleSSVMessage.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode Message from SSVMessage to aggregate commits")
	}
	message.FullData = fullData
	data, err := message.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode Message with FullData to aggregate commits")
	}
	sampleSSVMessage.Data = data

	// Create aggregated commit
	ret := &types.SignedSSVMessage{
		OperatorID: signers,
		Signature:  signatures,
		SSVMessage: sampleSSVMessage,
	}

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
func CreateCommit(state *State, config IConfig, root [32]byte) (*Message, error) {
	return &Message{
		MsgType:    CommitMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,

		Root: root,
	}, nil
}

func baseCommitValidation(
	config IConfig,
	signedCommit *types.SignedSSVMessage,
	height Height,
	operators []*types.Operator,
) error {

	// Decode
	message := &Message{}
	if err := message.Decode(signedCommit.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode Commit Message")
	}

	if message.MsgType != CommitMsgType {
		return errors.New("commit msg type is wrong")
	}
	if message.Height != height {
		return errors.New("wrong msg height")
	}

	if err := signedCommit.Validate(); err != nil {
		return errors.Wrap(err, "signed commit invalid")
	}

	return nil
}

func validateCommit(
	config IConfig,
	signedCommit *types.SignedSSVMessage,
	height Height,
	round Round,
	proposedSignedSSVMsg *types.SignedSSVMessage,
	operators []*types.Operator,
) error {
	if err := baseCommitValidation(config, signedCommit, height, operators); err != nil {
		return err
	}

	// Decode
	message := &Message{}
	if err := message.Decode(signedCommit.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode Commit Message")
	}

	if len(signedCommit.GetOperatorIDs()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if message.Round != round {
		return errors.New("wrong msg round")
	}

	// Decode
	proposedMessage := &Message{}
	if err := proposedMessage.Decode(proposedSignedSSVMsg.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode Proposal Message for commit validation")
	}

	if !bytes.Equal(proposedMessage.Root[:], message.Root[:]) {
		return errors.New("proposed data mistmatch")
	}

	return nil
}
