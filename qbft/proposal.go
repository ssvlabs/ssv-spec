package qbft

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// uponProposal process proposal message
// Assumes proposal message is valid!
func (i *Instance) uponProposal(msg *ProcessingMessage, proposeMsgContainer *MsgContainer) error {
	addedMsg, err := proposeMsgContainer.AddFirstMsgForSignerAndRound(msg)
	if err != nil {
		return errors.Wrap(err, "could not add proposal msg to container")
	}
	if !addedMsg {
		return nil // uponProposal was already called
	}

	newRound := msg.QBFTMessage.Round
	i.State.ProposalAcceptedForCurrentRound = msg

	// A future justified proposal should bump us into future round and reset timer
	if msg.QBFTMessage.Round > i.State.Round {
		i.config.GetTimer().TimeoutForRound(msg.QBFTMessage.Round)
	}
	i.State.Round = newRound

	// value root
	r, err := HashDataRoot(msg.SignedMessage.FullData)
	if err != nil {
		return errors.Wrap(err, "could not hash input data")
	}

	prepare, err := CreatePrepare(i.State, i.config, newRound, r)
	if err != nil {
		return errors.Wrap(err, "could not create prepare msg")
	}

	if err := i.Broadcast(prepare); err != nil {
		return errors.Wrap(err, "failed to broadcast prepare message")
	}

	return nil
}

func isValidProposal(
	state *State,
	config IConfig,
	msg *ProcessingMessage,
	valCheck ProposedValueCheckF,
) error {

	if msg.QBFTMessage.MsgType != ProposalMsgType {
		return errors.New("msg type is not proposal")
	}
	if msg.QBFTMessage.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(msg.SignedMessage.OperatorIDs) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if !msg.SignedMessage.CheckSignersInCommittee(state.CommitteeMember.Committee) {
		return errors.New("signer not in committee")
	}

	if !msg.SignedMessage.MatchedSigners([]types.OperatorID{proposer(state, config, msg.QBFTMessage.Round)}) {
		return errors.New("proposal leader invalid")
	}

	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "proposal invalid")
	}

	// verify full data integrity
	r, err := HashDataRoot(msg.SignedMessage.FullData)
	if err != nil {
		return errors.Wrap(err, "could not hash input data")
	}
	if !bytes.Equal(msg.QBFTMessage.Root[:], r[:]) {
		return errors.New("H(data) != root")
	}

	// get justifications
	roundChangeJustificationSignedMessages, _ := msg.QBFTMessage.GetRoundChangeJustifications() // no need to check error, checked on signedProposal.Validate()
	prepareJustificationSignedMessages, _ := msg.QBFTMessage.GetPrepareJustifications()         // no need to check error, checked on signedProposal.Validate()

	roundChangeJustification := make([]*ProcessingMessage, 0)
	for _, rcSignedMessage := range roundChangeJustificationSignedMessages {
		rc, err := NewProcessingMessage(rcSignedMessage)
		if err != nil {
			return errors.Wrap(err, "could not create ProcessingMessage from round change justification")
		}
		roundChangeJustification = append(roundChangeJustification, rc)
	}
	prepareJustification := make([]*ProcessingMessage, 0)
	for _, prepareSignedMessage := range prepareJustificationSignedMessages {
		msg, err := NewProcessingMessage(prepareSignedMessage)
		if err != nil {
			return errors.Wrap(err, "could not create ProcessingMessage from prepare justification")
		}
		prepareJustification = append(prepareJustification, msg)
	}

	if err := isProposalJustification(
		state,
		config,
		roundChangeJustification,
		prepareJustification,
		state.Height,
		msg.QBFTMessage.Round,
		msg.SignedMessage.FullData,
		valCheck,
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}

	if (state.ProposalAcceptedForCurrentRound == nil && msg.QBFTMessage.Round == state.Round) ||
		msg.QBFTMessage.Round > state.Round {
		return nil
	}
	return errors.New("proposal is not valid with current state")
}

// isProposalJustification returns nil if the proposal and round change messages are valid and justify a proposal message for the provided round, value and leader
func isProposalJustification(
	state *State,
	config IConfig,
	roundChangeMsgs []*ProcessingMessage,
	prepareMsgs []*ProcessingMessage,
	height Height,
	round Round,
	fullData []byte,
	valCheck ProposedValueCheckF,
) error {
	if err := valCheck(fullData); err != nil {
		return errors.Wrap(err, "proposal fullData invalid")
	}

	if round == FirstRound {
		return nil
	} else {
		// check all round changes are valid for height and round
		// no quorum, duplicate signers,  invalid still has quorum, invalid no quorum
		// prepared
		for _, rc := range roundChangeMsgs {
			if err := validRoundChangeForDataVerifySignature(state, config, rc, height, round, fullData); err != nil {
				return errors.Wrap(err, "change round msg not valid")
			}
		}

		// check there is a quorum
		if !HasQuorum(state.CommitteeMember, roundChangeMsgs) {
			return errors.New("change round has no quorum")
		}

		// previouslyPreparedF returns true if any on the round change messages have a prepared round and fullData
		previouslyPrepared, err := func(rcMsgs []*ProcessingMessage) (bool, error) {
			for _, rc := range rcMsgs {

				if rc.QBFTMessage.RoundChangePrepared() {
					return true, nil
				}
			}
			return false, nil
		}(roundChangeMsgs)
		if err != nil {
			return errors.Wrap(err, "could not calculate if previously prepared")
		}

		if !previouslyPrepared {
			return nil
		} else {

			// check prepare quorum
			if !HasQuorum(state.CommitteeMember, prepareMsgs) {
				return errors.New("prepares has no quorum")
			}

			// get a round change data for which there is a justification for the highest previously prepared round
			rcMsg, err := highestPrepared(roundChangeMsgs)
			if err != nil {
				return errors.Wrap(err, "could not get highest prepared")
			}
			if rcMsg == nil {
				return errors.New("no highest prepared")
			}

			// proposed fullData must equal highest prepared fullData
			r, err := HashDataRoot(fullData)
			if err != nil {
				return errors.Wrap(err, "could not hash input data")
			}
			if !bytes.Equal(r[:], rcMsg.QBFTMessage.Root[:]) {
				return errors.New("proposed data doesn't match highest prepared")
			}

			// validate each prepare message against the highest previously prepared fullData and round
			for _, pm := range prepareMsgs {
				if err := validSignedPrepareForHeightRoundAndRootVerifySignature(
					config,
					pm,
					height,
					rcMsg.QBFTMessage.DataRound,
					rcMsg.QBFTMessage.Root,
					state.CommitteeMember.Committee,
				); err != nil {
					return errors.New("signed prepare not valid")
				}
			}
			return nil
		}
	}
}

func proposer(state *State, config IConfig, round Round) types.OperatorID {
	// TODO - https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/29ae5a44551466453a84d4d17b9e083ecf189d97/dafny/spec/L1/node_auxiliary_functions.dfy#L304-L323
	return config.GetProposerF()(state, round)
}

// CreateProposal
/**
  	Proposal(
                        signProposal(
                            UnsignedProposal(
                                |current.blockchain|,
                                newRound,
                                digest(block)),
                            current.id),
                        block,
                        extractSignedRoundChanges(roundChanges),
                        extractSignedPrepares(prepares));
*/
func CreateProposal(state *State, config IConfig, fullData []byte, roundChanges, prepares []*ProcessingMessage) (*types.SignedSSVMessage, error) {
	r, err := HashDataRoot(fullData)
	if err != nil {
		return nil, errors.Wrap(err, "could not hash input data")
	}

	roundChangeSignedMessages := make([]*types.SignedSSVMessage, 0)
	for _, msg := range roundChanges {
		roundChangeSignedMessages = append(roundChangeSignedMessages, msg.SignedMessage)
	}
	prepareSignedMessages := make([]*types.SignedSSVMessage, 0)
	for _, msg := range prepares {
		prepareSignedMessages = append(prepareSignedMessages, msg.SignedMessage)
	}

	roundChangesData, err := MarshalJustifications(roundChangeSignedMessages)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal justifications")
	}
	preparesData, err := MarshalJustifications(prepareSignedMessages)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal justifications")
	}

	msg := &Message{
		MsgType:    ProposalMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,

		Root:                     r,
		RoundChangeJustification: roundChangesData,
		PrepareJustification:     preparesData,
	}

	signedMsg, err := Sign(msg, state.CommitteeMember.OperatorID, config.GetOperatorSigner())
	if err != nil {
		return nil, errors.Wrap(err, "could not create proposal message")
	}
	signedMsg.FullData = fullData
	return signedMsg, nil
}
