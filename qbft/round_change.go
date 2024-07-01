package qbft

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// uponRoundChange process round change messages.
// Assumes round change message is valid!
func (i *Instance) uponRoundChange(
	instanceStartValue []byte,
	msg *ProcessingMessage,
	roundChangeMsgContainer *MsgContainer,
	valCheck ProposedValueCheckF,
) error {

	hasQuorumBefore := HasQuorum(i.State.CommitteeMember, roundChangeMsgContainer.MessagesForRound(msg.QBFTMessage.Round))
	// Currently, even if we have a quorum of round change messages, we update the container
	addedMsg, err := roundChangeMsgContainer.AddFirstMsgForSignerAndRound(msg)
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
		msg,
		roundChangeMsgContainer,
		valCheck)
	if err != nil {
		return errors.Wrap(err, "could not get proposal justification for leading round")
	}

	if justifiedRoundChangeMsg != nil {

		roundChangeJustificationSignedMessages, _ := justifiedRoundChangeMsg.QBFTMessage.GetRoundChangeJustifications() // no need to check error, check on isValidRoundChange

		roundChangeJustification := make([]*ProcessingMessage, 0)
		for _, rcSignedMessage := range roundChangeJustificationSignedMessages {
			rc, err := NewProcessingMessage(rcSignedMessage)
			if err != nil {
				return errors.Wrap(err, "could not create ProcessingMessage from round change justification")
			}
			roundChangeJustification = append(roundChangeJustification, rc)
		}

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

func hasReceivedPartialQuorum(state *State, roundChangeMsgContainer *MsgContainer) (bool, []*ProcessingMessage) {
	all := roundChangeMsgContainer.AllMessages()

	rc := make([]*ProcessingMessage, 0)
	for _, msg := range all {

		if msg.QBFTMessage.Round > state.Round {
			rc = append(rc, msg)
		}
	}

	return HasPartialQuorum(state.CommitteeMember, rc), rc
}

// hasReceivedProposalJustificationForLeadingRound returns
// if first round or not received round change msgs with prepare justification - returns first rc msg in container and value to propose
// if received round change msgs with prepare justification - returns the highest prepare justification round change msg and value to propose
// (all the above considering the operator is a leader for the round
func hasReceivedProposalJustificationForLeadingRound(
	state *State,
	config IConfig,
	instanceStartValue []byte,
	roundChangeMessage *ProcessingMessage,
	roundChangeMsgContainer *MsgContainer,
	valCheck ProposedValueCheckF,
) (*ProcessingMessage, []byte, error) {

	roundChanges := roundChangeMsgContainer.MessagesForRound(roundChangeMessage.QBFTMessage.Round)

	// optimization, if no round change quorum can return false
	if !HasQuorum(state.CommitteeMember, roundChanges) {
		return nil, nil, nil
	}

	// Important!
	// We iterate on all round chance msgs for liveliness in case the last round change msg is malicious.
	for _, containerRoundChangeMessage := range roundChanges {

		// Chose proposal value.
		// If justifiedRoundChangeMsg has no prepare justification chose state value
		// If justifiedRoundChangeMsg has prepare justification chose prepared value
		valueToPropose := instanceStartValue
		if containerRoundChangeMessage.QBFTMessage.RoundChangePrepared() {
			valueToPropose = containerRoundChangeMessage.SignedMessage.FullData
		}

		roundChangeSignedMessagesJustification, _ := containerRoundChangeMessage.QBFTMessage.GetRoundChangeJustifications() // no need to check error, checked on isValidRoundChange

		roundChangeJustification := make([]*ProcessingMessage, 0)
		for _, signedMessage := range roundChangeSignedMessagesJustification {
			msg, err := NewProcessingMessage(signedMessage)
			if err != nil {
				return nil, nil, errors.Wrap(err, "could not create ProcessingMessage from round change justification")
			}
			roundChangeJustification = append(roundChangeJustification, msg)
		}

		if isProposalJustificationForLeadingRound(
			state,
			config,
			containerRoundChangeMessage,
			roundChanges,
			roundChangeJustification,
			valueToPropose,
			valCheck,
			roundChangeMessage.QBFTMessage.Round,
		) == nil {
			// not returning error, no need to
			return containerRoundChangeMessage, valueToPropose, nil
		}
	}
	return nil, nil, nil
}

// isProposalJustificationForLeadingRound - returns nil if we have a quorum of round change msgs and highest justified value for leading round
func isProposalJustificationForLeadingRound(
	state *State,
	config IConfig,
	roundChangeMsg *ProcessingMessage,
	roundChanges []*ProcessingMessage,
	roundChangeJustifications []*ProcessingMessage,
	value []byte,
	valCheck ProposedValueCheckF,
	newRound Round,
) error {

	if err := isReceivedProposalJustification(
		state,
		config,
		roundChanges,
		roundChangeJustifications,
		roundChangeMsg.QBFTMessage.Round,
		value,
		valCheck); err != nil {
		return err
	}

	if proposer(state, config, roundChangeMsg.QBFTMessage.Round) != state.CommitteeMember.OperatorID {
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
	roundChanges, prepares []*ProcessingMessage,
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

func validRoundChangeForDataIgnoreSignature(
	state *State,
	config IConfig,
	msg *ProcessingMessage,
	height Height,
	round Round,
	fullData []byte,
) error {

	if msg.QBFTMessage.MsgType != RoundChangeMsgType {
		return errors.New("round change msg type is wrong")
	}
	if msg.QBFTMessage.Height != height {
		return errors.New("wrong msg height")
	}
	if msg.QBFTMessage.Round != round {
		return errors.New("wrong msg round")
	}
	if len(msg.SignedMessage.OperatorIDs) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "roundChange invalid")
	}

	if !msg.SignedMessage.CheckSignersInCommittee(state.CommitteeMember.Committee) {
		return errors.New("signer not in committee")
	}

	// Addition to formal spec
	// We add this extra tests on the msg itself to filter round change msgs with invalid justifications, before they are inserted into msg containers
	if msg.QBFTMessage.RoundChangePrepared() {
		r, err := HashDataRoot(fullData)
		if err != nil {
			return errors.Wrap(err, "could not hash input data")
		}

		// validate prepare message justifications
		prepareSignedMsgs, _ := msg.QBFTMessage.GetRoundChangeJustifications() // no need to check error, checked on msg.QBFTMessage.Validate()

		prepareMsgs := make([]*ProcessingMessage, 0)
		for _, signedMessage := range prepareSignedMsgs {
			msg, err := NewProcessingMessage(signedMessage)
			if err != nil {
				return errors.Wrap(err, "could not create ProcessingMessage from prepare message in round change justification")
			}
			prepareMsgs = append(prepareMsgs, msg)
		}

		for _, pm := range prepareMsgs {
			if err := validSignedPrepareForHeightRoundAndRootVerifySignature(
				config,
				pm,
				state.Height,
				msg.QBFTMessage.DataRound,
				msg.QBFTMessage.Root,
				state.CommitteeMember.Committee); err != nil {
				return errors.Wrap(err, "round change justification invalid")
			}
		}

		if !bytes.Equal(r[:], msg.QBFTMessage.Root[:]) {
			return errors.New("H(data) != root")
		}

		if !HasQuorum(state.CommitteeMember, prepareMsgs) {
			return errors.New("no justifications quorum")
		}

		if msg.QBFTMessage.DataRound > round {
			return errors.New("prepared round > round")
		}

		return nil
	}

	return nil
}

func validRoundChangeForDataVerifySignature(
	state *State,
	config IConfig,
	msg *ProcessingMessage,
	height Height,
	round Round,
	fullData []byte,
) error {

	if err := validRoundChangeForDataIgnoreSignature(state, config, msg, height, round, fullData); err != nil {
		return err
	}

	// Verify signature
	if err := config.GetSignatureVerifier().Verify(msg.SignedMessage, state.CommitteeMember.Committee); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	return nil
}

// highestPrepared returns a round change message with the highest prepared round, returns nil if none found
func highestPrepared(roundChanges []*ProcessingMessage) (*ProcessingMessage, error) {
	var ret *ProcessingMessage
	var highestPreparedRound Round
	for _, rc := range roundChanges {

		if !rc.QBFTMessage.RoundChangePrepared() {
			continue
		}

		if ret == nil {
			ret = rc
			highestPreparedRound = rc.QBFTMessage.DataRound
		} else {
			if highestPreparedRound < rc.QBFTMessage.DataRound {
				ret = rc
				highestPreparedRound = rc.QBFTMessage.DataRound
			}
		}
	}
	return ret, nil
}

// returns the min round number out of the signed round change messages and the current round
func minRound(roundChangeMsgs []*ProcessingMessage) Round {
	ret := NoRound
	for _, msg := range roundChangeMsgs {
		if ret == NoRound || ret > msg.QBFTMessage.Round {
			ret = msg.QBFTMessage.Round
		}
	}
	return ret
}

func getRoundChangeData(state *State, config IConfig, instanceStartValue []byte) (Round, [32]byte, []byte, []*ProcessingMessage, error) {
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
func CreateRoundChange(state *State, config IConfig, newRound Round, instanceStartValue []byte) (*types.SignedSSVMessage, error) {
	round, root, fullData, justifications, err := getRoundChangeData(state, config, instanceStartValue)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate round change data")
	}

	justificationsSignedMessages := make([]*types.SignedSSVMessage, 0)
	for _, msg := range justifications {
		justificationsSignedMessages = append(justificationsSignedMessages, msg.SignedMessage)
	}

	justificationsData, err := MarshalJustifications(justificationsSignedMessages)
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
	}
	signedMsg, err := Sign(msg, state.CommitteeMember.OperatorID, config.GetOperatorSigner())
	if err != nil {
		return nil, errors.Wrap(err, "could not sign round change message")
	}
	signedMsg.FullData = fullData
	return signedMsg, nil
}
