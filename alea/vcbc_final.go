package alea

import (
	"bytes"
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponVCBCFinal(signedMessage *SignedMessage) error {
	if i.verbose {
		fmt.Println("uponVCBCFinal")
	}
	// get Data
	vcbcFinalData, err := signedMessage.Message.GetVCBCFinalData()
	if err != nil {
		errors.Wrap(err, "uponVCBCFinal: could not get vcbcFinalData data from signedMessage")
	}

	// check if it has the message locally. If not, returns (since it can't validate the hash)
	if !i.State.VCBCState.hasM(vcbcFinalData.Author, vcbcFinalData.Priority) {
		if i.verbose {
			fmt.Println("\tDidn't have the message locally")
		}
		return nil
	}

	proposals := i.State.VCBCState.getM(vcbcFinalData.Author, vcbcFinalData.Priority)

	// get hash
	localHash, err := GetProposalsHash(proposals)
	if err != nil {
		errors.Wrap(err, "uponVCBCFinal: could not get hash of local proposals")
	}
	if i.verbose {
		fmt.Println("\tgot hash")
	}

	// compare hash with reiceved one
	if !bytes.Equal(vcbcFinalData.Hash, localHash) {
		if i.verbose {
			fmt.Println("\tdifferent hash, quiting.")
		}
		return nil
	}

	// check if already has local aggregated signature. If so, returns (since it alreasy has and delivered the proposals).
	if i.State.VCBCState.hasU(vcbcFinalData.Author, vcbcFinalData.Priority) {
		if i.verbose {
			fmt.Println("\talready has proof, quiting.")
		}
		return nil
	}

	// store proof
	i.State.VCBCState.setU(vcbcFinalData.Author, vcbcFinalData.Priority, vcbcFinalData.AggregatedMsg)

	// create VCBCDeliver message and broadcasts
	proposals = i.State.VCBCState.getM(vcbcFinalData.Author, vcbcFinalData.Priority)

	if i.verbose {
		fmt.Println("\tAdding to VCBC output.")
	}
	i.AddVCBCOutput(proposals, vcbcFinalData.Priority, vcbcFinalData.Author)
	if i.verbose {
		fmt.Println("\tnew queue for", vcbcFinalData.Author, " and priority", vcbcFinalData.Priority, ":", i.State.VCBCState.queues[vcbcFinalData.Author])
	}

	return nil
}

func isValidVCBCFinal(
	state *State,
	config IConfig,
	signedProposal *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedProposal.Message.MsgType != VCBCFinalMsgType {
		return errors.New("msg type is not VCBCFinalMsgType")
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

	VCBCFinalData, err := signedProposal.Message.GetVCBCFinalData()
	if err != nil {
		return errors.Wrap(err, "could not get VCBCFinalData data")
	}
	if err := VCBCFinalData.Validate(); err != nil {
		return errors.Wrap(err, "VCBCFinalData invalid")
	}

	// author
	author := VCBCFinalData.Author
	authorInCommittee := false
	for _, opID := range operators {
		if opID.OperatorID == author {
			authorInCommittee = true
		}
	}
	if !authorInCommittee {
		return errors.Wrap(err, "author (OperatorID) doesn't exist in Committee")
	}

	// priority & hash
	priority := VCBCFinalData.Priority
	if state.VCBCState.hasM(author, priority) {
		localHash, err := GetProposalsHash(state.VCBCState.getM(author, priority))
		if err != nil {
			return errors.Wrap(err, "could not get local hash")
		}
		if !bytes.Equal(localHash, VCBCFinalData.Hash) {
			return errors.Wrap(err, "existing (priority,author) proposals have different hash")
		}
	}

	// AggregatedMsg
	aggregatedMsg := VCBCFinalData.AggregatedMsg
	var signedAggregatedMessage *SignedMessage
	signedAggregatedMessage.Decode(aggregatedMsg)

	if err := signedAggregatedMessage.Signature.VerifyByOperators(signedAggregatedMessage, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "aggregatedMsg signature invalid")
	}
	if len(signedAggregatedMessage.GetSigners()) < int(state.Share.Quorum) {
		return errors.Wrap(err, "aggregatedMsg signers don't reach quorum")
	}

	return nil
}

func CreateVCBCFinal(state *State, config IConfig, hash []byte, priority Priority, aggregatedMsg []byte, author types.OperatorID) (*SignedMessage, error) {
	vcbcFinalData := &VCBCFinalData{
		Hash:          hash,
		Priority:      priority,
		AggregatedMsg: aggregatedMsg,
		Author:        author,
	}
	dataByts, err := vcbcFinalData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode vcbcFinalData")
	}
	msg := &Message{
		MsgType:    VCBCFinalMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing filler msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
