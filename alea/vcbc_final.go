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
	if bytes.Compare(vcbcFinalData.Hash, localHash) != 0 {
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
	i.State.VCBCState.setU(vcbcFinalData.Author, vcbcFinalData.Priority, vcbcFinalData.Proof)

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

func CreateVCBCFinal(state *State, config IConfig, hash []byte, priority Priority, proof types.Signature, author types.OperatorID) (*SignedMessage, error) {
	vcbcFinalData := &VCBCFinalData{
		Hash:     hash,
		Priority: priority,
		Proof:    proof,
		Author:   author,
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
