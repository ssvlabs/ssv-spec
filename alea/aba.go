package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)

func (i *Instance) StartABA(vote byte) byte {


	i.State.ABAState.Vin = vote

	initMsg, err := CreateABAInit(i.State, i.config, vote, i.State.ABAState.Round)
	if err != nil {
		errors.Wrap(err,"failed to create ABA Init message")
	}
	i.Broadcast(initMsg)

	<- i.State.ABAState.Terminate
	return i.State.ABAState.Vdecided
}




func (i *Instance) uponABA(signedABA *SignedMessage, abaMsgContainer *MsgContainer) error {
    
	fmt.Println("uponABA function")

	// // OperatorID who sent the message (and first signed it)
	// senderID := signedVCBC.GetSigners()[0]

	// fmt.Println("\tgot senderID:",senderID)

	// // initializes queue if it doesn't exist
	// if _, exists := i.State.queues[senderID]; !exists {
	// 	i.State.queues[senderID] = NewVCBCQueue()
	// }

	// // gets the sender's associated queue
	// queue := i.State.queues[senderID]
	// fmt.Println("\tgot senderID's queue:",queue)


	// // gets the VCBC data from the message
	// vcbcData, err := signedVCBC.Message.GetVCBCData()
	// if err != nil {
    // 	return errors.Wrap(err, "could not get VCBCData from message")
	// }
	// fmt.Println("\tgot vcbcData:",vcbcData)


	// // check if it was already delivered
	// if i.State.S.hasProposalList(vcbcData.ProposalData) {
	// 	fmt.Println("\tlist of proposals from VCBC already contained in S queue (already delivered)")
	// 	return nil
	// }
	// fmt.Println("\tlocal S hasn't yet this list of proposal:",vcbcData)


	// queue.Enqueue(vcbcData.ProposalData, vcbcData.Priority)
	// fmt.Println("\tenqueueing proposal list and priority")
	// fmt.Println("\tnew queue:",i.State.queues[senderID])

	return nil
}


func CreateABA(state *State, config IConfig, vote byte, round Round) (*SignedMessage, error) {
	abaData := &ABAData{
		Vote:			vote,
		Round:			round,					
	}
	dataByts, err := abaData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode aba data")
	}
	msg := &Message{
		MsgType:    ABAMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing aba msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
