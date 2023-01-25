package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)

func (i *Instance) uponVCBCReady(signedMessage *SignedMessage, msgContainer *MsgContainer) error {   
	fmt.Println("uponVCBCReady function")

	// get Data
	vcbcReadyData, err := signedMessage.Message.GetVCBCReadyData()
	if err != nil{
		errors.Wrap(err, "could not get vcbcReadyData data from signedMessage")
	}

	// before adding, check if it was already received
	msgAlreadyReceived, err := msgContainer.HasMsg(signedMessage)
	if err != nil{
		errors.Wrap(err, "could not check if message has already been received")
	}

	// Add message to container
    msgContainer.AddMsg(signedMessage)

	// get sender of the message

	// if message hasn't been received and this is the author of the VCBC -> aggregate signature
	if !msgAlreadyReceived && vcbcReadyData.Author == i.State.Share.OperatorID {

		// update W, the list of signedMessages to be aggregated later
		W := i.State.VCBCW[vcbcReadyData.Author][vcbcReadyData.Priority]
		W = append(W,signedMessage)
		
		// update counter associated with author and priority
		r := i.State.VCBCr[vcbcReadyData.Author][vcbcReadyData.Priority]
		r = r + 1
		
		// if reached quorum, aggregate signatures and broadcast FINAL message
		if r >= i.State.Share.Quorum {
			aggregatedMessage, err := aggregateMsgs(W)
			if err != nil {
				return errors.Wrap(err,"unable to aggregate messages to produce VCBCFinal")
			}
			i.State.VCBCu[vcbcReadyData.Author][vcbcReadyData.Priority] = aggregatedMessage.Signature

			vcbcFinalMsg, err := CreateVCBCFinal(i.State, i.config, vcbcReadyData.Hash, vcbcReadyData.Priority,  aggregatedMessage.Signature, vcbcReadyData.Author)
			if err != nil {
				return errors.Wrap(err, "failed to create VCBCReady message with proof")
			}
			i.Broadcast(vcbcFinalMsg)

			
		}
	}
	
	return nil
}


func aggregateMsgs(msgs []*SignedMessage) (*SignedMessage, error) {
	if len(msgs) == 0 {
		return nil, errors.New("can't aggregate zero msgs")
	}

	var ret *SignedMessage
	for _, m := range msgs {
		if ret == nil {
			ret = m.DeepCopy()
		} else {
			if err := ret.Aggregate(m); err != nil {
				return nil, errors.Wrap(err, "could not aggregate msg")
			}
		}
	}
	return ret, nil
}



func CreateVCBCReady(state *State, config IConfig, hash []byte, priority Priority,author types.OperatorID) (*SignedMessage, error) {
	vcbcReadyData := &VCBCReadyData{
		Hash:			hash,	
		Priority: 		priority,
		// Proof:			proof,
		Author:			author,	
	}
	dataByts, err := vcbcReadyData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode vcbcReadyData")
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
