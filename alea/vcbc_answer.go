package alea

import (
	"bytes"
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponVCBCAnswer(signedMessage *SignedMessage) error {

	if i.verbose {
		fmt.Println("uponVCBCAnswer")
	}

	// get Data
	vcbcAnswerData, err := signedMessage.Message.GetVCBCAnswerData()
	if err != nil {
		return errors.Wrap(err, "uponVCBCAnswer: could not get data from signedMessage")
	}

	// check if has local aggregated signature
	hasLocalSignature := i.State.VCBCState.hasU(vcbcAnswerData.Author, vcbcAnswerData.Priority)
	if hasLocalSignature {
		if i.verbose {
			fmt.Println("already had proof, returning")
		}
		return nil
	}

	if i.State.VCBCState.HasM(vcbcAnswerData.Author, vcbcAnswerData.Priority) {
		if !i.State.VCBCState.EqualM(vcbcAnswerData.Author, vcbcAnswerData.Priority, vcbcAnswerData.Proposals) {
			return errors.New("answer has different proposals than stores ones")
		}
	}

	// update local values
	i.State.VCBCState.SetU(vcbcAnswerData.Author, vcbcAnswerData.Priority, vcbcAnswerData.AggregatedMsg)
	i.State.VCBCState.setM(vcbcAnswerData.Author, vcbcAnswerData.Priority, vcbcAnswerData.Proposals)

	// add vcbc output
	i.AddVCBCOutput(vcbcAnswerData.Proposals, vcbcAnswerData.Priority, vcbcAnswerData.Author)

	return nil
}

func isValidVCBCAnswer(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != VCBCAnswerMsgType {
		return errors.New("msg type is not VCBCAnswerMsgType")
	}
	if signedMsg.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedMsg.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if err := signedMsg.Signature.VerifyByOperators(signedMsg, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	VCBCAnswerData, err := signedMsg.Message.GetVCBCAnswerData()
	if err != nil {
		return errors.Wrap(err, "could not get VCBCAnswerData data")
	}
	if err := VCBCAnswerData.Validate(); err != nil {
		return errors.Wrap(err, "VCBCAnswerData invalid")
	}

	// author
	author := VCBCAnswerData.Author
	authorInCommittee := false
	for _, opID := range operators {
		if opID.OperatorID == author {
			authorInCommittee = true
		}
	}
	if !authorInCommittee {
		return errors.New("author (OperatorID) doesn't exist in Committee")
	}

	// priority
	priority := VCBCAnswerData.Priority
	if state.VCBCState.HasM(author, priority) {
		if !state.VCBCState.EqualM(author, priority, VCBCAnswerData.Proposals) {
			return errors.Wrap(err, "existing (priority,author) with different proposals")
		}
	}

	// AggregatedMsg
	aggregatedMsg := VCBCAnswerData.AggregatedMsg
	signedAggregatedMessage := &SignedMessage{}
	signedAggregatedMessage.Decode(aggregatedMsg)

	if err := signedAggregatedMessage.Signature.VerifyByOperators(signedAggregatedMessage, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "aggregatedMsg signature invalid")
	}
	if len(signedAggregatedMessage.GetSigners()) < int(state.Share.Quorum) {
		return errors.New("aggregatedMsg signers don't reach quorum")
	}

	// AggregatedMsg data
	vcbcReadyData, err := signedAggregatedMessage.Message.GetVCBCReadyData()
	if err != nil {
		return errors.Wrap(err, "could not get VCBCReadyData from given aggregated message")
	}
	givenHash, err := GetProposalsHash(VCBCAnswerData.Proposals)
	if err != nil {
		return errors.Wrap(err, "could not get hash from given proposals")
	}
	if !bytes.Equal(givenHash, vcbcReadyData.Hash) {
		return errors.New("hash of proposals given doesn't match hash in the VCBCReadyData of the aggregated message")
	}
	if vcbcReadyData.Author != VCBCAnswerData.Author {
		return errors.New("author given doesn't match author in the VCBCReadyData of the aggregated message")
	}
	if vcbcReadyData.Priority != VCBCAnswerData.Priority {
		return errors.New("priority given doesn't match priority in the VCBCReadyData of the aggregated message")
	}

	return nil
}

func CreateVCBCAnswer(state *State, config IConfig, proposals []*ProposalData, priority Priority, aggregatedMsg []byte, author types.OperatorID) (*SignedMessage, error) {
	vcbcAnswerData := &VCBCAnswerData{
		Proposals:     proposals,
		Priority:      priority,
		AggregatedMsg: aggregatedMsg,
		Author:        author,
	}
	dataByts, err := vcbcAnswerData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCAnswer: could not encode vcbcAnswerData")
	}
	msg := &Message{
		MsgType:    VCBCAnswerMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCAnswer: failed signing filler msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
