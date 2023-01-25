package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)

func (i *Instance) uponFiller(signedFiller *SignedMessage, fillerMsgContainer *MsgContainer) error {   
	fmt.Println("uponFiller function")


	fillerData, err := signedFiller.Message.GetFillerData()
	if err != nil{
		errors.Wrap(err, "could not get filler data from signedFiller")
	}

	// Add message to container
    fillerMsgContainer.AddMsg(signedFiller)

	entries := fillerData.Entries
	priorities := fillerData.Priorities
	operatorID := fillerData.OperatorID

	queue := i.State.queues[operatorID]

	_, localPriority := queue.Peek()

	for idx, priority := range priorities {
		if priority > localPriority {
			queue.Enqueue(entries[idx],priority)
		}
	}
	

	return nil
}


func CreateFiller(state *State, config IConfig, entries [][]*ProposalData, priorities []Priority, operatorID types.OperatorID) (*SignedMessage, error) {
	fillerData := &FillerData{
		Entries:		entries,	
		Priorities: 	priorities,
		OperatorID:		operatorID,	
	}
	dataByts, err := fillerData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode filler data")
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
