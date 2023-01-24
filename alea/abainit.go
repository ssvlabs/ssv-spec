package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponABAInit(signedABAInit *SignedMessage, abaInitMsgContainer *MsgContainer) error {
	
	fmt.Println("uponABAInit function")

	abaInitData, err := signedABAInit.Message.GetABAInitData()
	if err != nil{
		errors.Wrap(err, "could not get abainitdata from signedABAInit")
	}

	// check if the round is the same as the current round
	// if initData.Round != i.State.ABARound {
	// 	return errors.Wrap(err,"received ABA init message for different round")
	// }

	// // check if the message is already in the container
	// if abaInitMsgContainer.Has(signedABAInit.Message.Identifier) {
	// 	return nil
	// }

	// add the message to the container
	i.State.ABAState.ABAInitContainer.AddMsg(signedABAInit)

	// Increase counter and
	// if weak support (partial quorum) is achieved for a value different from the input, broadcast a new init
	// if strong support (quorum) is achieved, add to values and broadcast AUX
	if abaInitData.Vote == 1 {
		i.State.ABAState.Init1Counter += 1
		if i.State.ABAState.Vin != 1 && i.State.Share.PartialQuorum == i.State.ABAState.Init1Counter{
			initMsg, err := CreateABAInit(i.State, i.config, abaInitData.Vote, abaInitData.Round)
			if err != nil {
				errors.Wrap(err,"failed to create ABA Init message after weak support")
			}
			i.Broadcast(initMsg)
		}
		if i.State.Share.Quorum == i.State.ABAState.Init1Counter {
			auxMsg, err := CreateABAAux(i.State, i.config, abaInitData.Vote, abaInitData.Round)
			if err != nil {
				errors.Wrap(err,"failed to create ABA Aux message after strong init support")
			}
			i.Broadcast(auxMsg)

			// initializes queue if it doesn't exist
			if _, exists := i.State.ABAState.Values[i.State.ABAState.Round]; !exists {
				i.State.ABAState.Values[i.State.ABAState.Round] = make([]byte,0)
			}
			i.State.ABAState.Values[i.State.ABAState.Round] = append(i.State.ABAState.Values[i.State.ABAState.Round],byte(1))
		}
	} else {
		i.State.ABAState.Init0Counter += 1
		if i.State.ABAState.Vin != 0 && i.State.Share.PartialQuorum == i.State.ABAState.Init0Counter{
			initMsg, err := CreateABAInit(i.State, i.config, abaInitData.Vote, abaInitData.Round)
			if err != nil {
				errors.Wrap(err,"failed to create ABA Init message after weak support")
			}
			i.Broadcast(initMsg)
		}
		if i.State.Share.Quorum == i.State.ABAState.Init0Counter {
			auxMsg, err := CreateABAAux(i.State, i.config, abaInitData.Vote, abaInitData.Round)
			if err != nil {
				errors.Wrap(err,"failed to create ABA Aux message after strong init support")
			}
			i.Broadcast(auxMsg)

			// initializes queue if it doesn't exist
			if _, exists := i.State.ABAState.Values[i.State.ABAState.Round]; !exists {
				i.State.ABAState.Values[i.State.ABAState.Round] = make([]byte,0)
			}
			i.State.ABAState.Values[i.State.ABAState.Round] = append(i.State.ABAState.Values[i.State.ABAState.Round],byte(0))
		}
	}
	return nil

}



	// // check if we have weak support
	// if abaInitMsgContainer.Size() >= (i.State.Share.Threshold + 1) {
	// 	// check if we have already broadcasted the message
	// 	if i.State.Init1Counter > 0 {
	// 		return nil
	// 	}

	// 	// broadcast the message
	// 	i.Broadcast(signedABAInit)
	// 	i.State.Init1Counter++
	// }

	// // check if we have strong support
	// if abaInitMsgContainer.Size() >= (2*i.State.Share.Threshold + 1) {
	// 	// check if we have already counted the votes
	// 	if i.State.Values[i.State.ABARound] != nil {
	// 		return nil
	

func CreateABAInit(state *State, config IConfig, vote byte, round Round) (*SignedMessage, error) {
	ABAInitData := &ABAInitData{
		Vote:		vote,
		Round:		round,					
	}
	dataByts, err := ABAInitData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode abainit data")
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
		return nil, errors.Wrap(err, "failed signing abainit msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
