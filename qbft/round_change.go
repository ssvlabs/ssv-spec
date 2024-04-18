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

	roundChangeMessage, err := DecodeMessage(signedRoundChange.SSVMessage.Data)
	if err != nil {
		return err
	}

	hasQuorumBefore := HasQuorum(i.State.Share, roundChangeMsgContainer.MessagesForRound(roundChangeMessage.
		Round))
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

	signedJustifiedRoundChangeMsg, valueToPropose, err := hasReceivedProposalJustificationForLeadingRound(
		i.State,
		i.config,
		instanceStartValue,
		signedRoundChange,
		roundChangeMsgContainer,
		valCheck)
	if err != nil {
		return errors.Wrap(err, "could not get proposal justification for leading round")
	}

	if signedJustifiedRoundChangeMsg != nil {

		justifiedRoundChangeMsg, err := DecodeMessage(signedJustifiedRoundChangeMsg.SSVMessage.Data)
		if err != nil {
			return err
		}
		roundChangeJustification, _ := justifiedRoundChangeMsg.GetRoundChangeJustifications() // no need to check error, check on isValidRoundChange

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
	all := roundChangeMsgContainer.AllMessages()

	rc := make([]*types.SignedSSVMessage, 0)
	for _, signedMsg := range all {

		msg, err := DecodeMessage(signedMsg.SSVMessage.Data)
		if err != nil {
			continue
		}

		if msg.Round > state.Round {
			rc = append(rc, signedMsg)
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

	roundChangeMessage, err := DecodeMessage(signedRoundChange.SSVMessage.Data)
	if err != nil {
		return nil, nil, err
	}

	roundChanges := roundChangeMsgContainer.MessagesForRound(roundChangeMessage.Round)

	// optimization, if no round change quorum can return false
	if !HasQuorum(state.Share, roundChanges) {
		return nil, nil, nil
	}

	// Important!
	// We iterate on all round chance msgs for liveliness in case the last round change msg is malicious.
	for _, containerRoundChangeSignedMessage := range roundChanges {

		containerRoundChangeMessage, err := DecodeMessage(containerRoundChangeSignedMessage.SSVMessage.Data)
		if err != nil {
			return nil, nil, err
		}

		// Chose proposal value.
		// If justifiedRoundChangeMsg has no prepare justification chose state value
		// If justifiedRoundChangeMsg has prepare justification chose prepared value
		valueToPropose := instanceStartValue
		if containerRoundChangeMessage.RoundChangePrepared() {
			valueToPropose = containerRoundChangeSignedMessage.FullData
		}

		roundChangeJustification, _ := containerRoundChangeMessage.GetRoundChangeJustifications() // no need to check error, checked on isValidRoundChange
		if isProposalJustificationForLeadingRound(
			state,
			config,
			containerRoundChangeSignedMessage,
			roundChanges,
			roundChangeJustification,
			valueToPropose,
			valCheck,
			roundChangeMessage.Round,
		) == nil {
			// not returning error, no need to
			return containerRoundChangeSignedMessage, valueToPropose, nil
		}
	}
	return nil, nil, nil
}

// isProposalJustificationForLeadingRound - returns nil if we have a quorum of round change msgs and highest justified value for leading round
func isProposalJustificationForLeadingRound(
	state *State,
	config IConfig,
	roundChangeSignedMsg *types.SignedSSVMessage,
	roundChanges []*types.SignedSSVMessage,
	roundChangeJustifications []*types.SignedSSVMessage,
	value []byte,
	valCheck ProposedValueCheckF,
	newRound Round,
) error {

	roundChangeMsg, err := DecodeMessage(roundChangeSignedMsg.SSVMessage.Data)
	if err != nil {
		return err
	}

	if err := isReceivedProposalJustification(
		state,
		config,
		roundChanges,
		roundChangeJustifications,
		roundChangeMsg.Round,
		value,
		valCheck); err != nil {
		return err
	}

	if proposer(state, config, roundChangeMsg.Round) != state.Share.OperatorID {
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

func validRoundChangeForDataIgnoreSignature(
	state *State,
	config IConfig,
	signedMsg *types.SignedSSVMessage,
	height Height,
	round Round,
	fullData []byte,
) error {

	msg, err := DecodeMessage(signedMsg.SSVMessage.Data)
	if err != nil {
		return err
	}

	if msg.MsgType != RoundChangeMsgType {
		return errors.New("round change msg type is wrong")
	}
	if msg.Height != height {
		return errors.New("wrong msg height")
	}
	if msg.Round != round {
		return errors.New("wrong msg round")
	}
	if len(signedMsg.GetOperatorIDs()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if err := signedMsg.Validate(); err != nil {
		return errors.Wrap(err, "roundChange invalid")
	}

	if !signedMsg.CheckSignersInCommittee(state.Share.Committee) {
		return errors.New("signer not in committee")
	}

	// Addition to formal spec
	// We add this extra tests on the msg itself to filter round change msgs with invalid justifications, before they are inserted into msg containers
	if msg.RoundChangePrepared() {
		r, err := HashDataRoot(fullData)
		if err != nil {
			return errors.Wrap(err, "could not hash input data")
		}

		// validate prepare message justifications
		prepareMsgs, _ := msg.GetRoundChangeJustifications() // no need to check error, checked on msg.Validate()
		for _, pm := range prepareMsgs {
			if err := validSignedPrepareForHeightRoundAndRootVerifySignature(
				config,
				pm,
				state.Height,
				msg.DataRound,
				msg.Root,
				state.Share.Committee); err != nil {
				return errors.Wrap(err, "round change justification invalid")
			}
		}

		if !bytes.Equal(r[:], msg.Root[:]) {
			return errors.New("H(data) != root")
		}

		if !HasQuorum(state.Share, prepareMsgs) {
			return errors.New("no justifications quorum")
		}

		if msg.DataRound > round {
			return errors.New("prepared round > round")
		}

		return nil
	}

	return nil
}

func validRoundChangeForDataVerifySignature(
	state *State,
	config IConfig,
	signedMsg *types.SignedSSVMessage,
	height Height,
	round Round,
	fullData []byte,
) error {

	if err := validRoundChangeForDataIgnoreSignature(state, config, signedMsg, height, round, fullData); err != nil {
		return err
	}

	// Verify signature
	if err := config.GetSignatureVerifier().Verify(signedMsg, state.Share.Committee); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	return nil
}

// highestPrepared returns a round change message with the highest prepared round, returns nil if none found
func highestPrepared(roundChanges []*types.SignedSSVMessage) (*types.SignedSSVMessage, error) {
	var ret *types.SignedSSVMessage
	var highestPreparedRound Round
	for _, rc := range roundChanges {

		msg, err := DecodeMessage(rc.SSVMessage.Data)
		if err != nil {
			continue
		}

		if !msg.RoundChangePrepared() {
			continue
		}

		if ret == nil {
			ret = rc
			highestPreparedRound = msg.DataRound
		} else {
			if highestPreparedRound < msg.DataRound {
				ret = rc
				highestPreparedRound = msg.DataRound
			}
		}
	}
	return ret, nil
}

// returns the min round number out of the signed round change messages and the current round
func minRound(roundChangeMsgs []*types.SignedSSVMessage) Round {
	ret := NoRound
	for _, signedMsg := range roundChangeMsgs {

		msg, err := DecodeMessage(signedMsg.SSVMessage.Data)
		if err != nil {
			continue
		}

		if ret == NoRound || ret > msg.Round {
			ret = msg.Round
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
func CreateRoundChange(state *State, config IConfig, newRound Round, instanceStartValue []byte) (*types.SignedSSVMessage, error) {
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
	}
	return MessageToSignedSSVMessageWithFullData(msg, state.Share.OperatorID, config.GetOperatorSigner(), fullData)
}
