package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// uponRoundChange process round change messages.
// Assumes round change message is valid!
func (i *Instance) uponRoundChange(
	signedRoundChange *SignedMessage,
	roundChangeMsgContainer *MsgContainer,
	valCheck ProposedValueCheckF,
) error {
	hasQuorumBefore := HasQuorum(i.State.Share, roundChangeMsgContainer.MessagesForRound(signedRoundChange.Message.
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

	justifiedRoundChangeMsg, valueToPropose, err := hasReceivedProposalJustificationForLeadingRound(
		i.State,
		i.config,
		i.CdFetcher,
		signedRoundChange,
		roundChangeMsgContainer,
		valCheck)
	if err != nil {
		return errors.Wrap(err, "could not get proposal justification for leading round")
	}
	if justifiedRoundChangeMsg != nil {
		roundChangeJustification, _ := justifiedRoundChangeMsg.Message.GetRoundChangeJustifications() // no need to check error, check on isValidRoundChange

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

		err := i.uponChangeRoundPartialQuorum(newRound)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Instance) uponChangeRoundPartialQuorum(newRound Round) error {
	i.State.Round = newRound
	i.State.ProposalAcceptedForCurrentRound = nil
	i.config.GetTimer().TimeoutForRound(i.State.Round)
	roundChange, err := CreateRoundChange(i.State, i.config, newRound)
	if err != nil {
		return errors.Wrap(err, "failed to create round change message")
	}
	if err := i.Broadcast(roundChange); err != nil {
		return errors.Wrap(err, "failed to broadcast round change message")
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
// if first round or not received round change msgs with prepare justification - returns first rc msg in container and value to propose
// if received round change msgs with prepare justification - returns the highest prepare justification round change msg and value to propose
// (all the above considering the operator is a leader for the round
func hasReceivedProposalJustificationForLeadingRound(
	state *State,
	config IConfig,
	cdFetcher *types.DataFetcher,
	signedRoundChange *SignedMessage,
	roundChangeMsgContainer *MsgContainer,
	valCheck ProposedValueCheckF,
) (*SignedMessage, []byte, error) {
	roundChanges := roundChangeMsgContainer.MessagesForRound(signedRoundChange.Message.Round)

	// optimization, if no round change quorum can return false
	if !HasQuorum(state.Share, roundChanges) {
		return nil, nil, nil
	}

	// Important!
	// We iterate on all round chance msgs for liveliness in case the last round change msg is malicious.
	for _, msg := range roundChanges {

		roundChangeJustification, _ := msg.Message.GetRoundChangeJustifications() // no need to check error, checked on isValidRoundChange
		if valueToPropose, err := isProposalJustificationForLeadingRound(
			state,
			config,
			msg,
			roundChanges,
			roundChangeJustification,
			cdFetcher,
			valCheck,
			signedRoundChange.Message.Round,
		); err == nil {
			// not returning error, no need to
			return msg, valueToPropose, nil
		}
	}
	return nil, nil, nil
}

// isProposalJustificationForLeadingRound - returns proposal value if we have a quorum of round change msgs and
// highest justified value for leading round
func isProposalJustificationForLeadingRound(
	state *State,
	config IConfig,
	roundChangeMsg *SignedMessage,
	roundChanges []*SignedMessage,
	roundChangeJustifications []*SignedMessage,
	cdFetcher *types.DataFetcher,
	valCheck ProposedValueCheckF,
	newRound Round,
) ([]byte, error) {
	if proposer(state, config, roundChangeMsg.Message.Round) != state.Share.OperatorID {
		return nil, errors.New("not proposer")
	}

	valueToPropose, err := chooseValueToPropose(roundChangeMsg, cdFetcher)
	if err != nil {
		return nil, err
	}

	if err := isReceivedProposalJustification(
		state,
		config,
		roundChanges,
		roundChangeJustifications,
		roundChangeMsg.Message.Round,
		valueToPropose,
		valCheck); err != nil {
		return nil, err
	}

	currentRoundProposal := state.ProposalAcceptedForCurrentRound == nil && state.Round == newRound
	futureRoundProposal := newRound > state.Round

	if !currentRoundProposal && !futureRoundProposal {
		return nil, errors.New("proposal round mismatch")
	}

	return valueToPropose, nil
}

// chooseValueToPropose
// If justifiedRoundChangeMsg has no prepare justification choose state value
// If justifiedRoundChangeMsg has prepare justification choose prepared value
func chooseValueToPropose(roundChangeMsg *SignedMessage, cdFetcher *types.DataFetcher) ([]byte, error) {
	var valueToPropose []byte
	if roundChangeMsg.Message.RoundChangePrepared() {
		valueToPropose = roundChangeMsg.FullData
	} else {
		cd, err := cdFetcher.GetConsensusData()
		if err != nil {
			return nil, err
		}
		valueToPropose = cd
	}
	return valueToPropose, nil
}

// isReceivedProposalJustification - returns nil if we have a quorum of round change msgs and highest justified value
func isReceivedProposalJustification(
	state *State,
	config IConfig,
	roundChanges, prepares []*SignedMessage,
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
	signedMsg *SignedMessage,
	height Height,
	round Round,
	fullData []byte,
) error {
	if signedMsg.Message.MsgType != RoundChangeMsgType {
		return errors.New("round change msg type is wrong")
	}
	if signedMsg.Message.Height != height {
		return errors.New("wrong msg height")
	}
	if signedMsg.Message.Round != round {
		return errors.New("wrong msg round")
	}
	if len(signedMsg.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if err := signedMsg.Signature.VerifyByOperators(signedMsg, config.GetSignatureDomainType(), types.QBFTSignatureType, state.Share.Committee); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	if err := signedMsg.Message.Validate(); err != nil {
		return errors.Wrap(err, "roundChange invalid")
	}

	// Addition to formal spec
	// We add this extra tests on the msg itself to filter round change msgs with invalid justifications, before they are inserted into msg containers
	if signedMsg.Message.RoundChangePrepared() {
		r, err := HashDataRoot(fullData)
		if err != nil {
			return errors.Wrap(err, "could not hash input data")
		}

		// validate prepare message justifications
		prepareMsgs, _ := signedMsg.Message.GetRoundChangeJustifications() // no need to check error, checked on signedMsg.Message.Validate()
		for _, pm := range prepareMsgs {
			if err := validSignedPrepareForHeightRoundAndRoot(
				config,
				pm,
				state.Height,
				signedMsg.Message.DataRound,
				signedMsg.Message.Root,
				state.Share.Committee); err != nil {
				return errors.Wrap(err, "round change justification invalid")
			}
		}

		if !bytes.Equal(r[:], signedMsg.Message.Root[:]) {
			return errors.New("H(data) != root")
		}

		if !HasQuorum(state.Share, prepareMsgs) {
			return errors.New("no justifications quorum")
		}

		if signedMsg.Message.DataRound > round {
			return errors.New("prepared round > round")
		}

		return nil
	}
	return nil
}

// highestPrepared returns a round change message with the highest prepared round, returns nil if none found
func highestPrepared(roundChanges []*SignedMessage) (*SignedMessage, error) {
	var ret *SignedMessage
	for _, rc := range roundChanges {
		if !rc.Message.RoundChangePrepared() {
			continue
		}

		if ret == nil {
			ret = rc
		} else {
			if ret.Message.DataRound < rc.Message.DataRound {
				ret = rc
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

func getRoundChangeData(state *State, config IConfig) (Round, [32]byte, []byte, []*SignedMessage, error) {
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
func CreateRoundChange(state *State, config IConfig, newRound Round) (*SignedMessage, error) {
	round, root, fullData, justifications, err := getRoundChangeData(state, config)
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
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing prepare msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   *msg,

		FullData: fullData,
	}
	return signedMsg, nil
}
