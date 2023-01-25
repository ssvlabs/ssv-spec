package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponFillGap(signedFillGap *SignedMessage, fillgapMsgContainer *MsgContainer) error {   
	fmt.Println("uponFillGap function")

	// get data
	fillGapData, err := signedFillGap.Message.GetFillGapData()
	if err != nil{
		errors.Wrap(err, "could not get fillgap data from signedFillGap")
	}

	// Add message to container
    fillgapMsgContainer.AddMsg(signedFillGap)

	// get structure values
	operatorID := fillGapData.OperatorID
	priorityAsked := fillGapData.Priority

	// get the desired queue
	queue := i.State.queues[operatorID]
	// get highest local priority
	_, priority := queue.PeekLast()
	
	// if has more entries than the asker (sender of the message), sends FILLER message with local entries
	if priority >= priorityAsked {
		// init values, priority list
		returnValues := make([][]*ProposalData,0)
		returnPriorities := make([]Priority,0)

		// get local values and priorities
		values := queue.GetValues()
		priorities := queue.GetPriorities()

		// for each, test if priority if above and, if so, adds to the FILLER list
		for idx,priority := range priorities {
			if priority >= priorityAsked {
				returnValues = append(returnValues,values[idx])
				returnPriorities = append(returnPriorities,priority)
			}
		}

		// sends FILLER message
		fillerMsg, err := CreateFiller(i.State, i.config, returnValues, returnPriorities, operatorID)
		if err != nil {
			errors.Wrap(err,"failed to create Filler message")
		}
		// FIX ME : send only to sender of fillGap msg
		i.Broadcast(fillerMsg)
	}

	return nil
}

func CreateFillGap(state *State, config IConfig, operatorID types.OperatorID, priority Priority) (*SignedMessage, error) {
	fillgapData := &FillGapData{
		OperatorID:		operatorID,
		Priority:			priority,					
	}
	dataByts, err := fillgapData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode fillgap data")
	}
	msg := &Message{
		MsgType:    FillGapMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing fillgap msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
