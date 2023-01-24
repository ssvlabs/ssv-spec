package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponFillGap(signedFillGap *SignedMessage, fillgapMsgContainer *MsgContainer) error {   
	fmt.Println("uponFillGap function")
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
