package qbft

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// uponProposal process proposal message
// Assumes proposal message is valid!
func (i *Instance) uponProposal(signedProposal *types.SignedSSVMessage, proposeMsgContainer *MsgContainer) error {
	addedMsg, err := proposeMsgContainer.AddFirstMsgForSignerAndRound(signedProposal)
	if err != nil {
		return errors.Wrap(err, "could not add proposal msg to container")
	}
	if !addedMsg {
		return nil // uponProposal was already called
	}

	msg, err := DecodeMessage(signedProposal.SSVMessage.Data)
	if err != nil {
		return err
	}

	newRound := msg.Round
	i.State.ProposalAcceptedForCurrentRound = signedProposal

	// A future justified proposal should bump us into future round and reset timer
	if msg.Round > i.State.Round {
		i.config.GetTimer().TimeoutForRound(msg.Round)
	}
	i.State.Round = newRound

	// value root
	r, err := HashDataRoot(signedProposal.FullData)
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
	signedProposal *types.SignedSSVMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {

	msg, err := DecodeMessage(signedProposal.SSVMessage.Data)
	if err != nil {
		return err
	}

	if msg.MsgType != ProposalMsgType {
		return errors.New("msg type is not proposal")
	}
	if msg.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedProposal.GetOperatorIDs()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if !signedProposal.CheckSignersInCommittee(state.Share.Committee) {
		return errors.New("signer not in committee")
	}

	if !signedProposal.MatchedSigners([]types.OperatorID{proposer(state, config, msg.Round)}) {
		return errors.New("proposal leader invalid")
	}

	if err := signedProposal.Validate(); err != nil {
		return errors.Wrap(err, "proposal invalid")
	}

	// verify full data integrity
	r, err := HashDataRoot(signedProposal.FullData)
	if err != nil {
		return errors.Wrap(err, "could not hash input data")
	}
	if !bytes.Equal(msg.Root[:], r[:]) {
		return errors.New("H(data) != root")
	}

	// get justifications
	roundChangeJustification, _ := msg.GetRoundChangeJustifications() // no need to check error, checked on signedProposal.Validate()
	prepareJustification, _ := msg.GetPrepareJustifications()         // no need to check error, checked on signedProposal.Validate()

	if err := isProposalJustification(
		state,
		config,
		roundChangeJustification,
		prepareJustification,
		state.Height,
		msg.Round,
		signedProposal.FullData,
		valCheck,
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}

	if (state.ProposalAcceptedForCurrentRound == nil && msg.Round == state.Round) ||
		msg.Round > state.Round {
		return nil
	}
	return errors.New("proposal is not valid with current state")
}

// isProposalJustification returns nil if the proposal and round change messages are valid and justify a proposal message for the provided round, value and leader
func isProposalJustification(
	state *State,
	config IConfig,
	roundChangeMsgs []*types.SignedSSVMessage,
	prepareMsgs []*types.SignedSSVMessage,
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
		if !HasQuorum(state.Share, roundChangeMsgs) {
			return errors.New("change round has no quorum")
		}

		// previouslyPreparedF returns true if any on the round change messages have a prepared round and fullData
		previouslyPrepared, err := func(rcMsgs []*types.SignedSSVMessage) (bool, error) {
			for _, rc := range rcMsgs {

				msg, err := DecodeMessage(rc.SSVMessage.Data)
				if err != nil {
					continue
				}

				if msg.RoundChangePrepared() {
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
			if !HasQuorum(state.Share, prepareMsgs) {
				return errors.New("prepares has no quorum")
			}

			// get a round change data for which there is a justification for the highest previously prepared round
			rcSignedMsg, err := highestPrepared(roundChangeMsgs)
			if err != nil {
				return errors.Wrap(err, "could not get highest prepared")
			}
			if rcSignedMsg == nil {
				return errors.New("no highest prepared")
			}

			rcMsg, err := DecodeMessage(rcSignedMsg.SSVMessage.Data)
			if err != nil {
				return errors.Wrap(err, "highest prepared can't be decoded to Message")
			}

			// proposed fullData must equal highest prepared fullData
			r, err := HashDataRoot(fullData)
			if err != nil {
				return errors.Wrap(err, "could not hash input data")
			}
			if !bytes.Equal(r[:], rcMsg.Root[:]) {
				return errors.New("proposed data doesn't match highest prepared")
			}

			// validate each prepare message against the highest previously prepared fullData and round
			for _, pm := range prepareMsgs {
				if err := validSignedPrepareForHeightRoundAndRootVerifySignature(
					config,
					pm,
					height,
					rcMsg.DataRound,
					rcMsg.Root,
					state.Share.Committee,
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
func CreateProposal(state *State, config IConfig, fullData []byte, roundChanges, prepares []*types.SignedSSVMessage) (*types.SignedSSVMessage, error) {
	r, err := HashDataRoot(fullData)
	if err != nil {
		return nil, errors.Wrap(err, "could not hash input data")
	}

	roundChangesData, err := MarshalJustifications(roundChanges)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal justifications")
	}
	preparesData, err := MarshalJustifications(prepares)
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

	return MessageToSignedSSVMessageWithFullData(msg, state.Share.OperatorID, config.GetOperatorSigner(), fullData)
}
