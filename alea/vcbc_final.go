package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
	"bytes"
)

func (i *Instance) uponVCBCFinal(signedMessage *SignedMessage, msgContainer *MsgContainer) error {   
	fmt.Println("uponVCBCFinal function")

	// get Data
	vcbcFinalData, err := signedMessage.Message.GetVCBCFinalData()
	if err != nil{
		errors.Wrap(err, "could not get vcbcFinalData data from signedMessage")
	}

	// Add message to container
    msgContainer.AddMsg(signedMessage)

	localHash, err := GetProposalsHash(i.State.VCBCm[vcbcFinalData.Author][vcbcFinalData.Priority])
	if err != nil{
		errors.Wrap(err, "could not get hash of local proposals")
	}
	if bytes.Compare(vcbcFinalData.Hash,localHash) != 0 {
		return nil
	}
	// check if already has local aggregated signature
	if _, exists := i.State.VCBCu[vcbcFinalData.Author][vcbcFinalData.Priority]; exists {
		return nil
	}

	i.State.VCBCu[vcbcFinalData.Author][vcbcFinalData.Priority] = vcbcFinalData.Proof

	// create VCBC message and broadcasts
	proposals := i.State.VCBCm[vcbcFinalData.Author][vcbcFinalData.Priority]
	msgToBroadcast, err := CreateVCBC(i.State, i.config, proposals, vcbcFinalData.Priority)
	fmt.Println("\tcreated VCBC message to broadcast")
	if err != nil {
		return errors.Wrap(err, "failed to create VCBC message")
	}
   	i.Broadcast(msgToBroadcast)
	
	return nil
}


func CreateVCBCFinal(state *State, config IConfig, hash []byte, priority Priority, proof types.Signature,author types.OperatorID) (*SignedMessage, error) {
	vcbcFinalData := &VCBCFinalData{
		Hash:			hash,	
		Priority: 		priority,
		Proof:			proof,
		Author:			author,	
	}
	dataByts, err := vcbcFinalData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode vcbcFinalData")
	}
	msg := &Message{
		MsgType:    FillerMsgType,
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
