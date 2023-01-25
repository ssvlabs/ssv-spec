package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponFillGap(signedFillGap *SignedMessage, fillgapMsgContainer *MsgContainer) error {   
	fmt.Println("uponFillGap function")


	fillGapData, err := signedFillGap.Message.GetFillGapData()
	if err != nil{
		errors.Wrap(err, "could not get fillgap data from signedFillGap")
	}

	// Add message to container
    fillgapMsgContainer.AddMsg(signedFillGap)

	operatorID := fillGapData.OperatorID
	priorityAsked := fillGapData.Priority

	queue := i.State.queues[operatorID]
	_, priority := queue.Peek()
	
	if priority >= priorityAsked {
		returnValues := make([][]*ProposalData,0)
		returnPriorities := make([]Priority,0)
		values := queue.GetValues()
		priorities := queue.GetPriorities()
		for idx,priority := range priorities {
			if priority >= priorityAsked {
				returnValues = append(returnValues,values[idx])
				returnPriorities = append(returnPriorities,priority)
			}
		}

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
