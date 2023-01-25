package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponABAAux(signedABAAux *SignedMessage, abaAuxMsgContainer *MsgContainer) error {
	
	fmt.Println("uponABAAux function")

	// get message Data
	ABAAuxData, err := signedABAAux.Message.GetABAAuxData()
	if err != nil{
		errors.Wrap(err, "could not get ABAAuxData from signedABAAux")
	}

	// add the message to the containers
	i.State.ABAState.ABAAuxContainer.AddMsg(signedABAAux)
	abaAuxMsgContainer.AddMsg(signedABAAux)


	// increment counter
	if ABAAuxData.Vote == 1 {
		i.State.ABAState.Aux1Counter += 1
	} else {
		i.State.ABAState.Aux0Counter += 1
	}

	// if received 2f+1 AUX messages, try to send CONF
	if (i.State.ABAState.Aux1Counter + i.State.ABAState.Aux0Counter) >= i.State.Share.Quorum && !i.State.ABAState.SentConf {
		// produce list of values received by AUX
		auxVals := make([]byte,0)
		if i.State.ABAState.Aux1Counter > 0 {
			auxVals = append(auxVals,byte(1))
		}
		if i.State.ABAState.Aux0Counter > 0 {
			auxVals = append(auxVals,byte(0))
		}

		// determine how many values received by AUX are in the values list
		equalValues := 0
		for _, auxVal := range auxVals {
			for _, val := range i.State.ABAState.Values[ABAAuxData.Round] {
				if auxVal == val {
					equalValues += 1
				}
			}
		}

		// if every value exists in the values list (is contained), broadcast CONF with round values
		if equalValues == len(auxVals) {
			// broadcast CONF message
			confMsg, err := CreateABAConf(i.State, i.config, i.State.ABAState.Values[ABAAuxData.Round], ABAAuxData.Round)
			if err != nil {
				errors.Wrap(err,"failed to create ABA Conf message after strong support")
			}
			i.Broadcast(confMsg)

			// update sent flag
			i.State.ABAState.SentConf = true
		}

	}
	
	return nil
}



func CreateABAAux(state *State, config IConfig, vote byte, round Round) (*SignedMessage, error) {
	ABAAuxData := &ABAAuxData{
		Vote:		vote,
		Round:		round,					
	}
	dataByts, err := ABAAuxData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode abaaux data")
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
		return nil, errors.Wrap(err, "failed signing abaaux msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
