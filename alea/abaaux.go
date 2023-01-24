package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponABAAux(signedABAAux *SignedMessage, abaAuxMsgContainer *MsgContainer) error {
	
	fmt.Println("uponABAAux function")

	ABAAuxData, err := signedABAAux.Message.GetABAAuxData()
	if err != nil{
		errors.Wrap(err, "could not get ABAAuxData from signedABAAux")
	}

	// add the message to the container
	i.State.ABAState.ABAAuxContainer.AddMsg(signedABAAux)

	if ABAAuxData.Vote == 1 {
		i.State.ABAState.Aux1Counter += 1

		if i.State.ABAState.Aux1Counter >= i.State.Share.Quorum {
			contains := false
			for _,vote := range i.State.ABAState.Values[ABAAuxData.Round] {
				if vote == ABAAuxData.Vote {
					contains = true
				}
			}
			if contains {
				confMsg, err := CreateABAConf(i.State, i.config, i.State.ABAState.Values[ABAAuxData.Round], ABAAuxData.Round)
				if err != nil {
					errors.Wrap(err,"failed to create ABA Conf message after strong support")
				}
				i.Broadcast(confMsg)
			}
		}

	} else {
		i.State.ABAState.Aux0Counter += 1

		if i.State.ABAState.Aux1Counter >= i.State.Share.Quorum {
			contains := false
			for _,vote := range i.State.ABAState.Values[ABAAuxData.Round] {
				if vote == ABAAuxData.Vote {
					contains = true
				}
			}
			if contains {
				confMsg, err := CreateABAConf(i.State, i.config, i.State.ABAState.Values[ABAAuxData.Round], ABAAuxData.Round)
				if err != nil {
					errors.Wrap(err,"failed to create ABA Conf message after strong support")
				}
				i.Broadcast(confMsg)
			}
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
