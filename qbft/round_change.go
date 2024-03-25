package qbft

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// uponRoundChange process round change messages.
// Assumes round change message is valid!
func (i *Instance) uponRoundChange(
	instanceStartValue []byte,
	signedRoundChange *types.SignedSSVMessage,
	roundChangeMsgContainer *MsgContainer,
	valCheck ProposedValueCheckF,
) error {

	// Decode
	rcMessage := &Message{}
	if err := rcMessage.Decode(signedRoundChange.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode Round Change Message")
	}

	hasQuorumBefore := HasQuorum(i.State.Share, roundChangeMsgContainer.MessagesForRound(rcMessage.Round))
	// Currently, even if we have a quorum of round change messages, we update the container
	addedMsg, err := roundChangeMsgContainer.AddFirstMsgForSignerAndRound(signedRoundChange)
	if err != nil {
		return errors.Wrap(err, "could not add round change msg to container")
	}
	if !addedMsg {
		return nil // message was already added from signer
	}

	if hasQuorumBefore {
		return nil // already changed round
	}

	justifiedRoundChangeMsg, valueToPropose, err := hasReceivedProposalJustificationForLeadingRound(
		i.State,
		i.config,
		instanceStartValue,
		signedRoundChange,
		roundChangeMsgContainer,
		valCheck)
	if err != nil {
		return errors.Wrap(err, "could not get proposal justification for leading round")
	}
	if justifiedRoundChangeMsg != nil {
		// Decode
		justifiedRCMessage := &Message{}
		if err := justifiedRCMessage.Decode(justifiedRoundChangeMsg.SSVMessage.Data); err != nil {
			return errors.Wrap(err, "could not decode justified Round Change Message")
		}
		roundChangeJustification, _ := justifiedRCMessage.GetRoundChangeJustifications() // no need to check error, check on isValidRoundChange

		proposal, err := CreateProposal(
			i.State,
			i.config,
			valueToPropose,
			roundChangeMsgContainer.MessagesForRound(i.State.Round), // TODO - might be optimized to include only necessary quorum
			roundChangeJustification,
		)
		if err != nil {
			return errors.Wrap(err, "failed to create proposal")
		}

		if err := i.Broadcast(proposal); err != nil {
			return errors.Wrap(err, "failed to broadcast proposal message")
		}
	} else if partialQuorum, rcs := hasReceivedPartialQuorum(i.State, roundChangeMsgContainer); partialQuorum {
		newRound := minRound(rcs)
		if newRound <= i.State.Round {
			return nil // no need to advance round
		}

		err := i.uponChangeRoundPartialQuorum(newRound, instanceStartValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Instance) uponChangeRoundPartialQuorum(newRound Round, instanceStartValue []byte) error {
	i.State.Round = newRound
	i.State.ProposalAcceptedForCurrentRound = nil
	i.config.GetTimer().TimeoutForRound(i.State.Round)
	roundChange, err := CreateRoundChange(i.State, i.config, newRound, instanceStartValue)
	if err != nil {
		return errors.Wrap(err, "failed to create round change message")
	}
	if err := i.Broadcast(roundChange); err != nil {
		return errors.Wrap(err, "failed to broadcast round change message")
	}
	return nil
}

func hasReceivedPartialQuorum(state *State, roundChangeMsgContainer *MsgContainer) (bool, []*types.SignedSSVMessage) {
	all := roundChangeMsgContainer.AllMessaged()

	rc := make([]*types.SignedSSVMessage, 0)
	for _, msg := range all {

		// Decode
		rcMessage := &Message{}
		if err := rcMessage.Decode(msg.SSVMessage.Data); err != nil {
			continue
		}

		if rcMessage.Round > state.Round {
			rc = append(rc, msg)
		}
	}

	return HasPartialQuorum(state.Share, rc), rc
}

// hasReceivedProposalJustificationForLeadingRound returns
// if first round or not received round change msgs with prepare justification - returns first rc msg in container and value to propose
// if received round change msgs with prepare justification - returns the highest prepare justification round change msg and value to propose
// (all the above considering the operator is a leader for the round
func hasReceivedProposalJustificationForLeadingRound(
	state *State,
	config IConfig,
	instanceStartValue []byte,
	signedRoundChange *types.SignedSSVMessage,
	roundChangeMsgContainer *MsgContainer,
	valCheck ProposedValueCheckF,
) (*types.SignedSSVMessage, []byte, error) {

	// Decode
	rcMessage := &Message{}
	if err := rcMessage.Decode(signedRoundChange.SSVMessage.Data); err != nil {
		return nil, nil, errors.Wrap(err, "could not decode RoundChange Message to check if has received proposal justification")
	}

	roundChanges := roundChangeMsgContainer.MessagesForRound(rcMessage.Round)

	// optimization, if no round change quorum can return false
	if !HasQuorum(state.Share, roundChanges) {
		return nil, nil, nil
	}

	// Important!
	// We iterate on all round chance msgs for liveliness in case the last round change msg is malicious.
	for _, msg := range roundChanges {

		// Decode
		rcMessageI := &Message{}
		if err := rcMessageI.Decode(msg.SSVMessage.Data); err != nil {
			return nil, nil, errors.Wrap(err, "could not decode stored RoundChange Message to check if has received proposal justification")
		}

		// Chose proposal value.
		// If justifiedRoundChangeMsg has no prepare justification chose state value
		// If justifiedRoundChangeMsg has prepare justification chose prepared value
		valueToPropose := instanceStartValue
		if rcMessageI.RoundChangePrepared() {
			valueToPropose = rcMessageI.FullData
		}

		roundChangeJustification, _ := rcMessageI.GetRoundChangeJustifications() // no need to check error, checked on isValidRoundChange
		if isProposalJustificationForLeadingRound(
			state,
			config,
			msg,
			roundChanges,
			roundChangeJustification,
			valueToPropose,
			valCheck,
			rcMessageI.Round,
		) == nil {
			// not returning error, no need to
			return msg, valueToPropose, nil
		}
	}
	return nil, nil, nil
}

// isProposalJustificationForLeadingRound - returns nil if we have a quorum of round change msgs and highest justified value for leading round
func isProposalJustificationForLeadingRound(
	state *State,
	config IConfig,
	roundChangeMsg *types.SignedSSVMessage,
	roundChanges []*types.SignedSSVMessage,
	roundChangeJustifications []*types.SignedSSVMessage,
	value []byte,
	valCheck ProposedValueCheckF,
	newRound Round,
) error {

	// Decode
	rcMessage := &Message{}
	if err := rcMessage.Decode(roundChangeMsg.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode RoundChange Message for proposal justification")
	}

	if err := isReceivedProposalJustification(
		state,
		config,
		roundChanges,
		roundChangeJustifications,
		rcMessage.Round,
		value,
		valCheck); err != nil {
		return err
	}

	if proposer(state, config, rcMessage.Round) != state.Share.OperatorID {
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
	roundChanges, prepares []*types.SignedSSVMessage,
	newRound Round,
	value []byte,
	valCheck ProposedValueCheckF,
) error {
	if err := isProposalJustification(
		state,
		config,
		roundChanges,
		prepares,
		state.Height,
		newRound,
		value,
		valCheck,
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}
	return nil
}

func validRoundChangeForData(
	state *State,
	config IConfig,
	signedMsg *types.SignedSSVMessage,
	height Height,
	round Round,
	fullData []byte,
	verifySignature bool,
) error {

	// Decode
	rcMessage := &Message{}
	if err := rcMessage.Decode(signedMsg.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode RoundChange Message to validate")
	}

	if rcMessage.MsgType != RoundChangeMsgType {
		return errors.New("round change msg type is wrong")
	}
	if rcMessage.Height != height {
		return errors.New("wrong msg height")
	}
	if rcMessage.Round != round {
		return errors.New("wrong msg round")
	}
	if len(signedMsg.GetOperatorIDs()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if verifySignature {
		if err := types.VerifySignedSSVMessage(signedMsg, state.Share.Committee); err != nil {
			return errors.Wrap(err, "msg signature invalid")
		}
	}

	if err := rcMessage.Validate(); err != nil {
		return errors.Wrap(err, "roundChange invalid")
	}

	// Addition to formal spec
	// We add this extra tests on the msg itself to filter round change msgs with invalid justifications, before they are inserted into msg containers
	if rcMessage.RoundChangePrepared() {
		r, err := HashDataRoot(fullData)
		if err != nil {
			return errors.Wrap(err, "could not hash input data")
		}

		// validate prepare message justifications
		prepareMsgs, _ := rcMessage.GetRoundChangeJustifications() // no need to check error, checked on signedMsg.Message.Validate()
		for _, pm := range prepareMsgs {
			if err := validSignedPrepareForHeightRoundAndRoot(
				config,
				pm,
				state.Height,
				rcMessage.DataRound,
				rcMessage.Root,
				state.Share.Committee,
				true); err != nil {
				return errors.Wrap(err, "round change justification invalid")
			}
		}

		if !bytes.Equal(r[:], rcMessage.Root[:]) {
			return errors.New("H(data) != root")
		}

		if !HasQuorum(state.Share, prepareMsgs) {
			return errors.New("no justifications quorum")
		}

		if rcMessage.DataRound > round {
			return errors.New("prepared round > round")
		}

		return nil
	}
	return nil
}

// highestPrepared returns a round change message with the highest prepared round, returns nil if none found
func highestPrepared(roundChanges []*types.SignedSSVMessage) (*types.SignedSSVMessage, error) {
	var ret *types.SignedSSVMessage
	highestRound := NoRound
	for _, rc := range roundChanges {
		// Decode
		rcMessage := &Message{}
		if err := rcMessage.Decode(rc.SSVMessage.Data); err != nil {
			return nil, errors.Wrap(err, "could not decode RoundChange Message to compute the highest prepared message")
		}

		if !rcMessage.RoundChangePrepared() {
			continue
		}

		if ret == nil {
			ret = rc
			highestRound = rcMessage.DataRound
		} else {
			if highestRound < rcMessage.DataRound {
				ret = rc
				highestRound = rcMessage.DataRound
			}
		}
	}
	return ret, nil
}

// returns the min round number out of the signed round change messages and the current round
func minRound(roundChangeMsgs []*types.SignedSSVMessage) Round {
	ret := NoRound
	for _, msg := range roundChangeMsgs {
		// Decode
		rcMessage := &Message{}
		if err := rcMessage.Decode(msg.SSVMessage.Data); err != nil {
			continue
		}
		if ret == NoRound || ret > rcMessage.Round {
			ret = rcMessage.Round
		}
	}
	return ret
}

func getRoundChangeData(state *State, config IConfig, instanceStartValue []byte) (Round, [32]byte, []byte, []*types.SignedSSVMessage, error) {
	if state.LastPreparedRound != NoRound && state.LastPreparedValue != nil {
		justifications, err := getRoundChangeJustification(state, config, state.PrepareContainer)
		if err != nil {
			return NoRound, [32]byte{}, nil, nil, errors.Wrap(err, "could not get round change justification")
		}

		r, err := HashDataRoot(state.LastPreparedValue)
		if err != nil {
			return NoRound, [32]byte{}, nil, nil, errors.Wrap(err, "could not hash input data")
		}

		return state.LastPreparedRound, r, state.LastPreparedValue, justifications, nil
	}
	return NoRound, [32]byte{}, nil, nil, nil
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
func CreateRoundChange(state *State, config IConfig, newRound Round, instanceStartValue []byte) (*Message, error) {
	round, root, fullData, justifications, err := getRoundChangeData(state, config, instanceStartValue)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate round change data")
	}

	justificationsData, err := MarshalJustifications(justifications)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal justifications")
	}
	msg := &Message{
		MsgType:    RoundChangeMsgType,
		Height:     state.Height,
		Round:      newRound,
		Identifier: state.ID,

		Root:                     root,
		DataRound:                round,
		RoundChangeJustification: justificationsData,
		FullData:                 fullData,
	}
	return msg, nil
}
