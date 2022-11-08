package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponRoundChange(signedRoundChange *SignedMessage) error {
	if err := validRoundChange(i.State, i.config, signedRoundChange, i.State.Height, signedRoundChange.Message.Round); err != nil {
		return errors.Wrap(err, "round change msg invalid")
	}

	addedMsg, err := i.State.RoundChangeContainer.AddFirstMsgForSignerAndRound(signedRoundChange)
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
		i.State.RoundChangeContainer,
		i.config.GetValueCheckF())
	if err != nil {
		return errors.Wrap(err, "could not get proposal justification for leading round")
	}
	if justifiedRoundChangeMsg != nil {
		// Chose proposal value.
		// If justifiedRoundChangeMsg has no prepare justification chose state value
		// If justifiedRoundChangeMsg has prepare justification chose prepared value
		valueToPropose := i.StartValue
		if justifiedRoundChangeMsg.Message.Prepared() {
			// TODO<olegshmuelov>: validate that justified round change msg holds the complete input data
			valueToPropose = &Data{
				Root:   justifiedRoundChangeMsg.Message.InputRoot,
				Source: justifiedRoundChangeMsg.InputSource,
			}
		}

		proposeMsg, err := CreateProposal(
			i.State,
			i.config,
			valueToPropose,
			i.State.RoundChangeContainer.MessagesForRound(i.State.Round), // TODO - might be optimized to include only necessary quorum
			justifiedRoundChangeMsg.RoundChangeJustifications,
		)
		if err != nil {
			return errors.Wrap(err, "failed to create proposal")
		}

		proposalEncoded, err := proposeMsg.Encode()
		if err != nil {
			return errors.Wrap(err, "could not encode round change message")
		}

		if err = i.Broadcast(proposalEncoded, types.ConsensusProposeMsgType); err != nil {
			return errors.Wrap(err, "failed to broadcast proposal message")
		}
	} else if partialQuorum, rcs := hasReceivedPartialQuorum(i.State, i.State.RoundChangeContainer); partialQuorum {
		newRound := minRound(rcs)
		if newRound <= i.State.Round {
			return nil // no need to advance round
		}

		i.State.Round = newRound
		// TODO - should we reset timeout here for the new round?
		i.State.ProposalAcceptedForCurrentRound = nil

		rcMsg, err := CreateRoundChange(i.State, i.config, newRound)
		if err != nil {
			return errors.Wrap(err, "failed to create round change message")
		}

		rcEncoded, err := rcMsg.Encode()
		if err != nil {
			return errors.Wrap(err, "could not encode round change message")
		}

		if err = i.Broadcast(rcEncoded, types.ConsensusRoundChangeMsgType); err != nil {
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

	// Important!
	// We iterate on all round change msgs for liveliness in case the last round change msg is malicious.
	for _, rc := range roundChanges {
		if isReceivedProposalJustificationForLeadingRound(
			state,
			config,
			rc,
			roundChanges,
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
	roundChanges []*SignedMessage,
	valCheck ProposedValueCheckF,
	newRound Round,
) error {
	inputData := &Data{
		Root:   roundChangeMsg.Message.InputRoot,
		Source: roundChangeMsg.InputSource,
	}
	if err := isReceivedProposalJustification(
		state,
		config,
		roundChanges,
		roundChangeMsg.RoundChangeJustifications,
		roundChangeMsg.Message.Round,
		inputData,
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
	roundChangeJustifications,
	prepareJustifications []*SignedMessage,
	newRound Round,
	inputData *Data,
	valCheck ProposedValueCheckF,
) error {
	if err := isProposalJustification(
		state,
		config,
		roundChangeJustifications,
		prepareJustifications,
		state.Height,
		newRound,
		inputData,
		valCheck,
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}
	return nil
}

func validRoundChange(state *State, config IConfig, signedMsg *SignedMessage, height Height, round Round) error {
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

	if signedMsg.Message.Prepared() {
		if len(signedMsg.RoundChangeJustifications) == 0 {
			return errors.New("round change justification invalid")
		}

		// Addition to formal spec
		// We add this extra tests on the msg itself to filter round change msgs with invalid justifications, before they are inserted into msg containers
		for _, rcj := range signedMsg.RoundChangeJustifications {
			if err := validSignedPrepareForHeightRoundAndValue(
				config,
				rcj,
				state.Height,
				signedMsg.Message.PreparedRound,
				signedMsg.Message.InputRoot[:],
				state.Share.Committee,
			); err != nil {
				return errors.Wrap(err, "round change justification invalid")
			}
		}

		if !HasQuorum(state.Share, signedMsg.RoundChangeJustifications) {
			return errors.New("no justifications quorum")
		}

		if signedMsg.Message.PreparedRound > round {
			return errors.New("prepared round > round")
		}
	}
	return nil
}

// highestPrepared returns a round change message with the highest prepared round, returns nil if none found
func highestPrepared(roundChangeJustifications []*SignedMessage) *SignedMessage {
	var ret *SignedMessage
	for _, rcj := range roundChangeJustifications {
		if !rcj.Message.Prepared() {
			continue
		}

		if ret == nil {
			ret = rcj
		} else {
			if ret.Message.PreparedRound < rcj.Message.PreparedRound {
				ret = rcj
			}
		}
	}
	return ret
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

func getRoundChangeData(state *State, config IConfig) ([]*SignedMessage, Round, *Data) {
	if state.LastPreparedRound != NoRound && state.LastPreparedValue != nil {
		justifications := getRoundChangeJustification(state, config, state.PrepareContainer)
		return justifications, state.LastPreparedRound, state.LastPreparedValue
	}
	return nil, NoRound, &Data{}
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
func CreateRoundChange(state *State, config IConfig, newRound Round) (*SignedMessage, error) {
	justifications, round, preparedValue := getRoundChangeData(state, config)

	msg := &Message{
		Height:        state.Height,
		Round:         newRound,
		InputRoot:     preparedValue.Root,
		PreparedRound: round,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing round change msg")
	}

	return &SignedMessage{
		Message:                   msg,
		Signature:                 sig,
		Signers:                   []types.OperatorID{state.Share.OperatorID},
		InputSource:               preparedValue.Source,
		RoundChangeJustifications: justifications,
	}, nil
}
