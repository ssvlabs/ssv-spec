package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponRoundChange(
	instanceStartValue []byte,
	signedRoundChange *SignedMessage,
	roundChangeMsgContainer *MsgContainer,
	valCheck ProposedValueCheckF,
) error {
	// TODO - Roberto comment: could happen we received a round change before we switched the round and this msg will be rejected (lost)
	if err := validRoundChange(i.State, i.config, signedRoundChange, i.State.Height, signedRoundChange.Message.Round); err != nil {
		return errors.Wrap(err, "round change msg invalid")
	}

	addedMsg, err := roundChangeMsgContainer.AddFirstMsgForSignerAndRound(signedRoundChange)
	if err != nil {
		return errors.Wrap(err, "could not add round change msg to container")
	}
	if !addedMsg {
		return nil // UponCommit was already called
	}

	justifiedRoundChangeMsg, err := hasReceivedProposalJustificationForLeadingRound(
		i.State,
		i.config,
		signedRoundChange,
		roundChangeMsgContainer,
		valCheck)
	if err != nil {
		return errors.Wrap(err, "could not get proposal justification for leading round")
	}
	if justifiedRoundChangeMsg != nil {
		//highestRCData, err := justifiedRoundChangeMsg.Message.GetRoundChangeData()
		//if err != nil {
		//	return errors.Wrap(err, "could not round change data from highestJustifiedRoundChangeMsg")
		//}

		// Chose proposal value.
		// If justifiedRoundChangeMsg has no prepare justification chose state value
		// If justifiedRoundChangeMsg has prepare justification chose prepared value
		valueToPropose := instanceStartValue

		//if highestRCData.Prepared() {
		//	valueToPropose = highestRCData.PreparedValue
		//}
		if justifiedRoundChangeMsg.Message.PreparedRound != NoRound || len(justifiedRoundChangeMsg.Message.Input) != 0 {
			valueToPropose = justifiedRoundChangeMsg.Message.Input
		}

		proposal, err := CreateProposal(
			i.State,
			i.config,
			valueToPropose,
			roundChangeMsgContainer.MessagesForRound(i.State.Round), // TODO - might be optimized to include only necessary quorum
			justifiedRoundChangeMsg.RoundChangeJustifications,
		)
		if err != nil {
			return errors.Wrap(err, "failed to create proposal")
		}

		proposalEncoded, err := proposal.Encode()
		if err != nil {
			return errors.Wrap(err, "could not encode proposal message")
		}

		msgID := types.PopulateMsgType(i.State.ID, types.ConsensusProposeMsgType)

		broadcastMsg := &types.Message{
			ID:   msgID,
			Data: proposalEncoded,
		}

		if err = i.Broadcast(broadcastMsg); err != nil {
			return errors.Wrap(err, "failed to broadcast proposal message")
		}
	} else if partialQuorum, rcs := hasReceivedPartialQuorum(i.State, roundChangeMsgContainer); partialQuorum {
		newRound := minRound(rcs)
		if newRound <= i.State.Round {
			return nil // no need to advance round
		}

		i.State.Round = newRound
		// TODO - should we reset timeout here for the new round?
		i.State.ProposalAcceptedForCurrentRound = nil

		roundChange, err := CreateRoundChange(i.State, i.config, newRound, instanceStartValue)
		if err != nil {
			return errors.Wrap(err, "failed to create round change message")
		}

		roundChangeEncoded, err := roundChange.Encode()
		if err != nil {
			return errors.Wrap(err, "could not encode round change message")
		}

		msgID := types.PopulateMsgType(i.State.ID, types.ConsensusRoundChangeMsgType)

		broadcastMsg := &types.Message{
			ID:   msgID,
			Data: roundChangeEncoded,
		}

		if err = i.Broadcast(broadcastMsg); err != nil {
			return errors.Wrap(err, "failed to broadcast round change message")
		}
	}
	return nil
}

func hasReceivedPartialQuorum(state *State, roundChangeMsgContainer *MsgContainer) (bool, []*SignedMessage) {
	all := roundChangeMsgContainer.AllMessaged()

	rc := make([]*SignedMessage, 0)
	for _, msg := range all {
		if msg.Message.Round > state.Round {
			rc = append(rc, msg)
		}
	}

	return HasPartialQuorum(state.Share, rc), rc
}

// hasReceivedProposalJustificationForLeadingRound returns
// if first round or not received round change msgs with prepare justification - returns first rc msg in container
// if received round change msgs with prepare justification - returns the highest prepare justification round change msg
// (all the above considering the operator is a leader for the round
func hasReceivedProposalJustificationForLeadingRound(
	state *State,
	config IConfig,
	signedRoundChange *SignedMessage,
	roundChangeMsgContainer *MsgContainer,
	valCheck ProposedValueCheckF,
) (*SignedMessage, error) {
	roundChanges := roundChangeMsgContainer.MessagesForRound(state.Round)

	// optimization, if no round change quorum can return false
	if !HasQuorum(state.Share, roundChanges) {
		return nil, nil
	}

	rcHeaders := make([]*SignedMessageHeader, 0)
	for _, rc := range roundChanges {
		rcHeader, err := rc.ToSignedMessageHeader()
		if err != nil {
			return nil, errors.Wrap(err, "could not convert signed msg to signed msg header")
		}
		rcHeaders = append(rcHeaders, rcHeader)
	}

	// Important!
	// We iterate on all round change msgs for liveliness in case the last round change msg is malicious.
	for _, rc := range roundChanges {
		if isReceivedProposalJustificationForLeadingRound(
			state,
			config,
			rc,
			rcHeaders,
			valCheck,
			signedRoundChange.Message.Round,
		) == nil {
			// not returning error, no need to
			return rc, nil
		}
	}
	return nil, nil
}

// isReceivedProposalJustificationForLeadingRound - returns nil if we have a quorum of round change msgs and highest justified value for leading round
func isReceivedProposalJustificationForLeadingRound(
	state *State,
	config IConfig,
	roundChangeMsg *SignedMessage,
	roundChanges []*SignedMessageHeader,
	valCheck ProposedValueCheckF,
	newRound Round,
) error {
	//rcData, err := roundChangeMsg.Message.GetRoundChangeData()
	//if err != nil {
	//	return errors.Wrap(err, "could not get round change data")
	//}

	if err := isReceivedProposalJustification(
		state,
		config,
		roundChanges,
		roundChangeMsg.RoundChangeJustifications,
		roundChangeMsg.Message.Round,
		roundChangeMsg.Message.Input,
		valCheck,
	); err != nil {
		return err
	}

	if proposer(state, config, roundChangeMsg.Message.Round) != state.Share.OperatorID {
		return errors.New("not proposer")
	}

	currentRoundProposal := state.ProposalAcceptedForCurrentRound == nil && state.Round == newRound
	futureRoundProposal := newRound > state.Round

	if !currentRoundProposal && !futureRoundProposal {
		return errors.New("proposal round mismatch")
	}

	return nil
}

// isReceivedProposalJustification - returns nil if we have a quorum of round change msgs and highest justified value
func isReceivedProposalJustification(
	state *State,
	config IConfig,
	roundChangeJustifications []*SignedMessageHeader,
	prepareJustifications []*SignedMessageHeader,
	newRound Round,
	value []byte,
	valCheck ProposedValueCheckF,
) error {
	if err := isProposalJustification(
		state,
		config,
		roundChangeJustifications,
		prepareJustifications,
		state.Height,
		newRound,
		value,
		valCheck,
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}
	return nil
}

func validRoundChange(state *State, config IConfig, signedMsg *SignedMessage, height Height, round Round) error {
	//if signedMsg.Message.MsgType != RoundChangeMsgType {
	//	return errors.New("round change msg type is wrong")
	//}
	if signedMsg.Message.Height != height {
		return errors.New("round change Height is wrong")
	}
	if signedMsg.Message.Round != round {
		return errors.New("msg round wrong")
	}
	if len(signedMsg.GetSigners()) != 1 {
		return errors.New("round change msg allows 1 signer")
	}
	if err := signedMsg.Signature.VerifyByOperators(signedMsg, config.GetSignatureDomainType(), types.QBFTSignatureType, state.Share.Committee); err != nil {
		return errors.Wrap(err, "round change msg signature invalid")
	}

	// TODO<olegshmuelov>: len(signedMsg.Message.Input) is always != 0 (this check already done in message validation),
	// so is it always prepared?
	// check if prepared
	if signedMsg.Message.PreparedRound != NoRound || len(signedMsg.Message.Input) != 0 {
		// this check already done in message validation
		//	if len(d.PreparedValue) == 0 {
		//		return errors.New("round change prepared value invalid")
		//	}
		if len(signedMsg.RoundChangeJustifications) == 0 {
			return errors.New("round change justification invalid")
		}
		// TODO - should next proposal data be equal to prepared value?

		// Addition to formal spec
		// We add this extra tests on the msg itself to filter round change msgs with invalid justifications, before they are inserted into msg containers
		for _, rcj := range signedMsg.RoundChangeJustifications {
			if err := validSignedPrepareHeaderForHeightRoundAndValue(
				config,
				rcj,
				state.Height,
				signedMsg.Message.PreparedRound,
				signedMsg.Message.Input,
				state.Share.Committee,
			); err != nil {
				return errors.Wrap(err, "round change justification invalid")
			}
		}

		if !HasQuorumHeaders(state.Share, signedMsg.RoundChangeJustifications) {
			return errors.New("no justifications quorum")
		}

		if signedMsg.Message.PreparedRound > round {
			return errors.New("prepared round > round")
		}
	}
	return nil
}

// TODO<olegshmuelov>: remove returning error if not needed
// highestPrepared returns a round change message with the highest prepared round, returns nil if none found
func highestPrepared(roundChangesJustifications []*SignedMessageHeader) (*SignedMessageHeader, error) {
	var ret *SignedMessageHeader
	for _, rcj := range roundChangesJustifications {
		//rcData, err := rc.Message.GetRoundChangeData()
		//if err != nil {
		//	return nil, errors.Wrap(err, "could not get round change data")
		//}
		//
		//if !rcData.Prepared() {
		//	continue
		//}

		// TODO<olegshmuelov> check for input root
		//if rcj.Message.PreparedRound == NoRound && len(d.PreparedValue) == 0 {
		if rcj.Message.PreparedRound == NoRound {
			continue
		}

		if ret == nil {
			ret = rcj
		} else {
			//retRCData, err := ret.Message.GetRoundChangeData()
			//if err != nil {
			//	return nil, errors.Wrap(err, "could not get round change data")
			//}
			//if retRCData.PreparedRound < rcData.PreparedRound {
			//	ret = rc
			//}
			if ret.Message.PreparedRound < rcj.Message.PreparedRound {
				ret = rcj
			}
		}
	}
	return ret, nil
}

// returns the min round number out of the signed round change messages and the current round
func minRound(roundChangeMsgs []*SignedMessage) Round {
	ret := NoRound
	for _, msg := range roundChangeMsgs {
		if ret == NoRound || ret > msg.Message.Round {
			ret = msg.Message.Round
		}
	}
	return ret
}

func getRoundChangeData(state *State, config IConfig) ([]*SignedMessageHeader, Round, []byte) {
	if state.LastPreparedRound != NoRound && state.LastPreparedValue != nil {
		justifications := getRoundChangeJustification(state, config, state.PrepareContainer)
		return justifications, state.LastPreparedRound, state.LastPreparedValue
	}
	return nil, NoRound, nil
}

// CreateRoundChange
/**
RoundChange(
           signRoundChange(
               UnsignedRoundChange(
                   |current.blockchain|,
                   newRound,
                   digestOptionalBlock(current.lastPreparedBlock),
                   current.lastPreparedRound),
           current.id),
           current.lastPreparedBlock,
           getRoundChangeJustification(current)
       )
*/
func CreateRoundChange(state *State, config IConfig, newRound Round, instanceStartValue []byte) (*SignedMessage, error) {
	justifications, round, preparedValue := getRoundChangeData(state, config)
	//rcData, err := getRoundChangeData(state, config, instanceStartValue)
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not generate round change data")
	//}
	//dataByts, err := rcData.Encode()
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not encode round change data")
	//}

	msg := &Message{
		Height:        state.Height,
		Round:         newRound,
		Input:         preparedValue,
		PreparedRound: round,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing prepare msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,

		RoundChangeJustifications: justifications,
	}
	return signedMsg, nil
}
