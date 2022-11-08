package qbft

import (
	"bytes"
	"github.com/pkg/errors"

	"github.com/bloxapp/ssv-spec/types"
)

func (i *Instance) uponProposal(signedProposal *SignedMessage) error {
	valCheck := i.config.GetValueCheckF()
	if err := isValidProposal(i.State, i.config, signedProposal, valCheck, i.State.Share.Committee); err != nil {
		return errors.Wrap(err, "proposal invalid")
	}

	addedMsg, err := i.State.ProposeContainer.AddFirstMsgForSignerAndRound(signedProposal)
	if err != nil {
		return errors.Wrap(err, "could not add proposal msg to container")
	}
	if !addedMsg {
		return nil // uponProposal was already called
	}

	newRound := signedProposal.Message.Round
	i.State.ProposalAcceptedForCurrentRound = signedProposal

	// A future justified proposal should bump us into future round and reset timer
	if signedProposal.Message.Round > i.State.Round {
		i.config.GetTimer().TimeoutForRound(signedProposal.Message.Round)
	}
	i.State.Round = newRound

	prepareMsg, err := CreatePrepare(i.State, i.config, newRound, signedProposal.Message.InputRoot)
	if err != nil {
		return errors.Wrap(err, "could not create prepare msg")
	}

	prepareEncoded, err := prepareMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode prepare message")
	}

	if err = i.Broadcast(prepareEncoded, types.ConsensusPrepareMsgType); err != nil {
		return errors.Wrap(err, "failed to broadcast prepare message")
	}

	return nil
}

func isValidProposal(
	state *State,
	config IConfig,
	signedProposal *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedProposal.Message.Height != state.Height {
		return errors.New("proposal Height is wrong")
	}
	if len(signedProposal.GetSigners()) != 1 {
		return errors.New("proposal msg allows 1 signer")
	}
	if err := signedProposal.Signature.VerifyByOperators(signedProposal, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "proposal msg signature invalid")
	}
	if !signedProposal.MatchedSigners([]types.OperatorID{proposer(state, config, signedProposal.Message.Round)}) {
		return errors.New("proposal leader invalid")
	}

	inputData := &Data{
		Root:   signedProposal.Message.InputRoot,
		Source: signedProposal.InputSource,
	}
	if err := isProposalJustification(
		state,
		config,
		signedProposal.RoundChangeJustifications,
		signedProposal.ProposalJustifications,
		state.Height,
		signedProposal.Message.Round,
		inputData,
		valCheck,
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}

	if (state.ProposalAcceptedForCurrentRound == nil && signedProposal.Message.Round == state.Round) ||
		signedProposal.Message.Round > state.Round {
		return nil
	}
	return errors.New("proposal is not valid with current state")
}

// isProposalJustification returns nil if the proposal and round change messages are valid and justify a proposal message for the provided round, value and leader
func isProposalJustification(
	state *State,
	config IConfig,
	roundChangeJustifications []*SignedMessage,
	prepareJustifications []*SignedMessage,
	height Height,
	round Round,
	inputData *Data,
	valCheck ProposedValueCheckF,
) error {
	if err := valCheck(inputData); err != nil {
		return errors.Wrap(err, "proposal value invalid")
	}

	if round == FirstRound {
		return nil
	} else {
		// check all round changes are valid for height and round
		// no quorum, duplicate signers, invalid still has quorum, invalid no quorum
		// prepared
		for _, rcj := range roundChangeJustifications {
			if err := validRoundChange(state, config, rcj, height, round); err != nil {
				return errors.Wrap(err, "change round msg not valid")
			}
		}

		// check there is a quorum
		if !HasQuorum(state.Share, roundChangeJustifications) {
			return errors.New("change round has no quorum")
		}

		// previouslyPreparedF returns true if any on the round change messages have a prepared round and value
		previouslyPrepared, err := func(rcJustifications []*SignedMessage) (bool, error) {
			for _, rc := range rcJustifications {
				if rc.Message.Prepared() {
					return true, nil
				}
			}
			return false, nil
		}(roundChangeJustifications)
		if err != nil {
			return errors.Wrap(err, "could not calculate if previously prepared")
		}

		if !previouslyPrepared {
			return nil
		} else {

			// check prepare quorum
			if !HasQuorum(state.Share, prepareJustifications) {
				return errors.New("prepares has no quorum")
			}

			// get a round change data for which there is a justification for the highest previously prepared round
			rch := highestPrepared(roundChangeJustifications)
			if rch == nil {
				return errors.New("no highest prepared")
			}

			// proposed value must equal highest prepared value
			if !bytes.Equal(inputData.Root[:], rch.Message.InputRoot[:]) {
				return errors.New("proposed data doesn't match highest prepared")
			}

			// validate each prepare message against the highest previously prepared value and round
			for _, pj := range prepareJustifications {
				if err := validSignedPrepareForHeightRoundAndValue(
					config,
					pj,
					height,
					rch.Message.PreparedRound,
					rch.Message.InputRoot[:],
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
func CreateProposal(
	state *State,
	config IConfig,
	value *Data,
	roundChanges,
	prepares []*SignedMessage,
) (*SignedMessage, error) {
	msg := &Message{
		Height:    state.Height,
		Round:     state.Round,
		InputRoot: value.Root,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing proposal msg")
	}

	proposeMsg := &SignedMessage{
		Message:                   msg,
		Signers:                   []types.OperatorID{state.Share.OperatorID},
		Signature:                 sig,
		InputSource:               value.Source,
		RoundChangeJustifications: roundChanges,
		ProposalJustifications:    prepares,
	}

	return proposeMsg, nil
}
