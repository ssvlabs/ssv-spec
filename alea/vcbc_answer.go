package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)

func (i *Instance) uponVCBCAnswer(signedMessage *SignedMessage, msgContainer *MsgContainer) error {   
	fmt.Println("uponVCBCAnswer function")

	// get Data
	vcbcAnswerData, err := signedMessage.Message.GetVCBCAnswerData()
	if err != nil{
		errors.Wrap(err, "could not get data from signedMessage")
	}

	// Add message to container
    msgContainer.AddMsg(signedMessage)
	fmt.Println("\tAdded message to container")

	// check if has local aggregated signature
	_, exists := i.State.VCBCu[vcbcAnswerData.Author][vcbcAnswerData.Priority]
	if exists {
		return nil
	}

	// update local values
	i.State.VCBCu[vcbcAnswerData.Author][vcbcAnswerData.Priority] = vcbcAnswerData.Proof
	i.State.VCBCm[vcbcAnswerData.Author][vcbcAnswerData.Priority] = vcbcAnswerData.Proposals
	

	// deliver: create VCBC message and broadcasts
	msgToBroadcast, err := CreateVCBC(i.State, i.config, vcbcAnswerData.Proposals, vcbcAnswerData.Priority)
	fmt.Println("\tcreated VCBC message to broadcast")
	if err != nil {
		return errors.Wrap(err, "failed to create VCBC message")
	}
   	i.Broadcast(msgToBroadcast)

	return nil
}


func CreateVCBCAnswer(state *State, config IConfig, proposals []*ProposalData, priority Priority, proof types.Signature, author types.OperatorID) (*SignedMessage, error) {
	vcbcAnswerData := &VCBCAnswerData{
		Proposals:		proposals,	
		Priority: 		priority,
		Proof:			proof,
		Author:			author,	
	}
	dataByts, err := vcbcAnswerData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode vcbcAnswerData")
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
