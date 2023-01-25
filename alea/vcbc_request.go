package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)

func (i *Instance) uponVCBCRequest(signedMessage *SignedMessage, msgContainer *MsgContainer) error {   
	fmt.Println("uponVCBCRequest function")

	// get Data
	vcbcRequestData, err := signedMessage.Message.GetVCBCRequestData()
	if err != nil{
		errors.Wrap(err, "could not get data from signedMessage")
	}

	// Add message to container
    msgContainer.AddMsg(signedMessage)
	fmt.Println("\tAdded message to container")

	// check if has local aggregated signature
	u, exists := i.State.VCBCu[vcbcRequestData.Author][vcbcRequestData.Priority]
	if !exists {
		return nil
	}

	proposals := i.State.VCBCm[vcbcRequestData.Author][vcbcRequestData.Priority]
	msgToBroadcast, err := CreateVCBCAnswer(i.State, i.config, proposals, vcbcRequestData.Priority, u, vcbcRequestData.Author)
	fmt.Println("\tcreated VCBCAnswer message to broadcast")
	if err != nil {
		return errors.Wrap(err, "failed to create VCBCAnswer message")
	}

	// FIX ME : send only to requester
   	i.Broadcast(msgToBroadcast)
	
	return nil
}


func CreateVCBCRequest(state *State, config IConfig, priority Priority, author types.OperatorID) (*SignedMessage, error) {
	vcbcRequestData := &VCBCRequestData{	
		Priority: 		priority,
		Author:			author,	
	}
	dataByts, err := vcbcRequestData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode vcbcRequestData")
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
