package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)

func (i *Instance) uponVCBCBroadcast(signedMessage *SignedMessage, msgContainer *MsgContainer) error {   
	fmt.Println("uponVCBCBroadcast function")
	

	// get Data
	vcbcBroadcastData, err := signedMessage.Message.GetVCBCBroadcastData()
	if err != nil{
		errors.Wrap(err, "could not get data from signedMessage")
	}

	// Add message to container
    msgContainer.AddMsg(signedMessage)
	fmt.Println("\tAdded message to container")

	// create VCBCSend message and broadcasts
	msgToBroadcast, err := CreateVCBCSend(i.State, i.config, vcbcBroadcastData.Proposals, vcbcBroadcastData.Priority, vcbcBroadcastData.Author)
	if err != nil {
		return errors.Wrap(err, "failed to create VCBCSend message")
	}
    i.Broadcast(msgToBroadcast)
	
	return nil
}


func CreateVCBCBroadcast(state *State, config IConfig, proposals []*ProposalData, priority Priority, author types.OperatorID) (*SignedMessage, error) {
	vcbcBroadcastData := &VCBCBroadcastData{
		Proposals:		proposals,	
		Priority: 		priority,
		Author:			author,	
	}
	dataByts, err := vcbcBroadcastData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode vcbcBroadcastData")
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
