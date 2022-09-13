package qbft

import (
	"github.com/pkg/errors"

	"github.com/bloxapp/ssv-spec/types"
)

func (i *Instance) uponProposal(signedProposal *SignedMessage, proposeMsgContainer *MsgContainer) error {
	valCheck := i.config.GetValueCheckF()
	if err := isValidProposal(i.State, i.config, signedProposal, valCheck, i.State.Share.Committee); err != nil {
		return errors.Wrap(err, "proposal invalid")
	}

	addedMsg, err := proposeMsgContainer.AddFirstMsgForSignerAndRound(signedProposal)
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

	//proposalData, err := signedProposal.Message.GetProposalData()
	//if err != nil {
	//	return errors.Wrap(err, "could not get proposal data")
	//}

	prepare, err := CreatePrepare(i.State, i.config, newRound, signedProposal.Message.Input)
	if err != nil {
		return errors.Wrap(err, "could not create prepare msg")
	}

	prepareEncoded, err := prepare.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode prepare message")
	}

	msgID := types.PopulateMsgType(i.State.ID, types.ConsensusPrepareMsgType)

	broadcastMsg := &types.Message{
		ID:   msgID,
		Data: prepareEncoded,
	}

	if err = i.Broadcast(broadcastMsg); err != nil {
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
	//if signedProposal.Message.MsgType != ProposalMsgType {
	//	return errors.New("msg type is not proposal")
	//}
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

	//proposalData, err := signedProposal.Message.GetProposalData()
	//if err != nil {
	//	return errors.Wrap(err, "could not get proposal data")
	//}
	//if err := proposalData.Validate(); err != nil {
	//	return errors.Wrap(err, "proposalData invalid")
	//}

	if err := isProposalJustification(
		state,
		config,
		signedProposal.RoundChangeJustifications,
		signedProposal.ProposalJustifications,
		//proposalData.RoundChangeJustification,
		//proposalData.PrepareJustification,
		state.Height,
		signedProposal.Message.Round,
		signedProposal.Message.Input,
		//proposalData.Data,
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
	roundChangeJustifications []*SignedMessageHeader,
	prepareJustifications []*SignedMessageHeader,
	height Height,
	round Round,
	value []byte,
	valCheck ProposedValueCheckF,
) error {
	if err := valCheck(value); err != nil {
		return errors.Wrap(err, "proposal value invalid")
	}

	if round == FirstRound {
		return nil
	} else {
		// check all round changes are valid for height and round
		// no quorum, duplicate signers, invalid still has quorum, invalid no quorum
		// prepared
		for _, rcj := range roundChangeJustifications {
			//if err := validRoundChange(state, config, rcj, height, round); err != nil {
			//	return errors.Wrap(err, "change round msg not valid")
			//}
			if err := validSignedPrepareHeaderForHeightRoundAndValue(
				config,
				rcj,
				state.Height,
				// TODO<olegshmuelov>: check for validity of the passes round and the value
				round,
				value,
				state.Share.Committee,
			); err != nil {
				return errors.Wrap(err, "round change justification invalid")
			}
		}

		// check there is a quorum
		if !HasQuorumHeaders(state.Share, roundChangeJustifications) {
			return errors.New("change round has not quorum")
		}

		// previouslyPreparedF returns true if any on the round change messages have a prepared round and value
		previouslyPrepared, err := func(rcJustifications []*SignedMessageHeader) (bool, error) {
			for _, rcj := range rcJustifications {
				//rcData, err := rcj.Message.GetRoundChangeData()
				//if err != nil {
				//	return false, errors.Wrap(err, "could not get round change data")
				//}
				//if rcData.Prepared() {
				//	return true, nil
				//}

				// TODO<olegshmuelov> check for input root
				//if rcj.Message.PreparedRound != NoRound || len(d.PreparedValue) != 0 {
				if rcj.Message.PreparedRound != NoRound {
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
			if !HasQuorumHeaders(state.Share, prepareJustifications) {
				return errors.New("prepares has no quorum")
			}

			// get a round change data for which there is a justification for the highest previously prepared round
			rcm, err := highestPrepared(roundChangeJustifications)
			if err != nil {
				return errors.Wrap(err, "could not get highest prepared")
			}
			if rcm == nil {
				return errors.New("no highest prepared")
			}
			//rcmData, err := rcm.Message.GetRoundChangeData()
			//if err != nil {
			//	return errors.Wrap(err, "could not get round change data")
			//}

			// TODO<olegshmuelov>: encode value to input root and compare?
			// proposed value must equal highest prepared value
			//if !bytes.Equal(value, rcm.Message.InputRoot) {
			//	return errors.New("proposed data doesn't match highest prepared")
			//}

			// validate each prepare message against the highest previously prepared value and round
			for _, pm := range prepareJustifications {
				if err := validSignedPrepareHeaderForHeightRoundAndValue(
					config,
					pm,
					height,
					rcm.Message.PreparedRound,
					// TODO<olegshmuelov>: finalize the comparison between input root in msgHeaders with inputValue
					// added [:] meanwhile
					rcm.Message.InputRoot[:],
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
	value []byte,
	roundChanges []*SignedMessage,
	prepares []*SignedMessageHeader,
) (*SignedMessage, error) {
	msg := &Message{
		Height: state.Height,
		Round:  state.Round,
		Input:  value,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing proposal msg")
	}

	proposalMsg := &SignedMessage{
		Message:                msg,
		Signers:                []types.OperatorID{state.Share.OperatorID},
		Signature:              sig,
		ProposalJustifications: prepares,
	}

	if len(roundChanges) > 0 {
		rcHeaders := make([]*SignedMessageHeader, 0)
		for _, rc := range roundChanges {
			rcHeader, err := rc.ToSignedMessageHeader()
			if err != nil {
				return nil, errors.Wrap(err, "could not convert signed msg to signed msg header")
			}
			rcHeaders = append(rcHeaders, rcHeader)
		}
		proposalMsg.RoundChangeJustifications = rcHeaders
	}

	//proposalData := &ProposalData{
	//	Data:                     value,
	//	RoundChangeJustification: roundChanges,
	//	PrepareJustification:     prepares,
	//}
	//dataByts, err := proposalData.Encode()

	//msg := &Message{
	//	MsgType:    ProposalMsgType,
	//	Height:     state.Height,
	//	Round:      state.Round,
	//	Identifier: state.ID,
	//	Data:       dataByts,
	//}
	//sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed signing prepare msg")
	//}
	//
	//signedMsg := &SignedMessage{
	//	Signature: sig,
	//	Signers:   []types.OperatorID{state.Share.OperatorID},
	//	Message:   msg,
	//}
	//return signedMsg, nil
	return proposalMsg, nil
}
