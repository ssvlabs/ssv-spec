package qbft

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// uponProposal process proposal message
// Assumes proposal message is valid!
func (i *Instance) uponProposal(signedProposal *SignedMessage, proposeMsgContainer *MsgContainer) error {
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
	signedProposal *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedProposal.Message.MsgType != ProposalMsgType {
		return errors.New("msg type is not proposal")
	}
	if signedProposal.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedProposal.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if err := signedProposal.Signature.VerifyByOperators(signedProposal, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}
	if !signedProposal.MatchedSigners([]types.OperatorID{proposer(state, config, signedProposal.Message.Round)}) {
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
	if !bytes.Equal(signedProposal.Message.Root[:], r[:]) {
		return errors.New("H(data) != root")
	}

	// get justifications
	roundChangeJustification, _ := signedProposal.Message.GetRoundChangeJustifications() // no need to check error, checked on signedProposal.Validate()
	proposalJustification, _ := signedProposal.Message.GetProposalJustifications()       // no need to check error, checked on signedProposal.Validate()

	if err := isProposalJustification(
		state,
		config,
		proposalJustification,
		roundChangeJustification,
		state.Height,
		signedProposal.Message.Round,
		signedProposal.FullData,
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
	roundChangeMsgs []*SignedMessage,
	prepareMsgs []*SignedMessage,
	height Height,
	round Round,
	fullData []byte,
	valCheck ProposedValueCheckF,
) error {
	// Check proposed data
	if err := valCheck(fullData); err != nil {
		return errors.Wrap(err, "proposal fullData invalid")
	}

	// If first round -> no validation to do
	if round == FirstRound {
		return nil
	}

	// Check unique signers of RC messages reach quorum
	if !HasQuorum(state.Share, roundChangeMsgs) {
		return errors.New("change round has no quorum")
	}

	// Validate each round change without validating the nested prepare messages
	for _, rc := range roundChangeMsgs {
		if err := validRoundChange(state, config, rc, height, round, fullData); err != nil {
			return errors.Wrap(err, "change round msg not valid")
		}
	}

	// Check if at least one RC message is prepared
	hasPreparedRoundChange := false
	for _, rc := range roundChangeMsgs {
		if rc.Message.RoundChangePrepared() {
			hasPreparedRoundChange = true
			break
		}
	}
	if !hasPreparedRoundChange {
		return nil
	}

	// If at least one Round-Change message is prepared,
	// we validate the proposed quorum of prepare messages against the highest prepared

	// check prepare quorum
	if !HasQuorum(state.Share, prepareMsgs) {
		return errors.New("prepares has no quorum")
	}

	// get a round change data for which there is a justification for the highest previously prepared round
	rcm, err := highestPrepared(roundChangeMsgs)
	if err != nil {
		return errors.Wrap(err, "could not get highest prepared")
	}
	if rcm == nil {
		return errors.New("no highest prepared")
	}

	checkHighest := func(rcm *SignedMessage, expectedRoot [32]byte) error {
		// proposed root must equal highest prepared root
		if !bytes.Equal(expectedRoot[:], rcm.Message.Root[:]) {
			return errors.New("proposed data doesn't match highest prepared")
		}

		// validate each prepare message against the highest previously prepared fullData and round
		for _, pm := range prepareMsgs {
			if err := validSignedPrepareForHeightRoundAndRoot(
				config,
				pm,
				height,
				rcm.Message.DataRound,
				rcm.Message.Root,
				state.Share.Committee,
			); err != nil {
				return errors.New("signed prepare not valid")
			}
		}
		return nil
	}

	highestPreparedRCs := highestPreparedSet(roundChangeMsgs)
	r, err := HashDataRoot(fullData)
	if err != nil {
		return errors.Wrap(err, "could not hash input data")
	}
	for _, rcm := range highestPreparedRCs {
		if checkHighest(rcm, r) == nil {
			return nil
		}
	}
	return errors.New("No highest prepared round-change matches prepared messages")
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
func CreateProposal(state *State, config IConfig, fullData []byte, roundChanges, prepares []*SignedMessage) (*SignedMessage, error) {
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
		RoundChangeJustification: preparesData,
		ProposalJustification:    roundChangesData,
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
