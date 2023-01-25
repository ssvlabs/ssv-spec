package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)

func (i *Instance) uponVCBCSend(signedMessage *SignedMessage, msgContainer *MsgContainer) error {   
	fmt.Println("uponVCBCSend function")

	// get Data
	vcbcSendData, err := signedMessage.Message.GetVCBCSendData()
	if err != nil{
		errors.Wrap(err, "could not get vcbcSendData data from signedMessage")
	}

	// before adding, check if it was already received
	msgAlreadyReceived, err := msgContainer.HasMsg(signedMessage)
	if err != nil{
		errors.Wrap(err, "could not check if message has already been received")
	}
	// Add message to container
    msgContainer.AddMsg(signedMessage)

	// get sender of the message
	senderID := signedMessage.GetSigners()[0]

	// if message hasn't been received and the Author of the VCBC is the same as the sender of the message -> sign and answer with READY
	if senderID == vcbcSendData.Author && !msgAlreadyReceived {
		
		hash, err := GetProposalsHash(vcbcSendData.Proposals)
		if err != nil{
			errors.Wrap(err, "could not get hash of proposals")
		}

		// create VCBCready message with nil proof
		// vcbcReadyMsg, err := CreateVCBCReady(i.State, i.config, hash, vcbcSendData.Priority, nil,vcbcSendData.Author)
		// if err != nil {
		// 	return errors.Wrap(err, "failed to create VCBCReady message with nil proof")
		// }

		// sign the VCBCReady message with nil proof
		// sig, err := i.config.GetSigner().SignRoot(vcbcReadyMsg, types.QBFTSignatureType, i.State.Share.SharePubKey)
		// if err != nil {
		// 	return errors.Wrap(err, "failed to sign root of vcbcReady msg")
		// }

		// create VCBCready message with proof
		vcbcReadyMsgWithSign, err := CreateVCBCReady(i.State, i.config, hash, vcbcSendData.Priority, vcbcSendData.Author)
		if err != nil {
			return errors.Wrap(err, "failed to create VCBCReady message with proof")
		}

		// FIX ME : send specifically to author
		i.Broadcast(vcbcReadyMsgWithSign)
	}
	
	return nil
}


func CreateVCBCSend(state *State, config IConfig, proposals []*ProposalData, priority Priority, author types.OperatorID) (*SignedMessage, error) {
	vcbcSendData := &VCBCSendData{
		Proposals:		proposals,	
		Priority: 		priority,
		Author:			author,	
	}
	dataByts, err := vcbcSendData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode vcbcSendData")
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
