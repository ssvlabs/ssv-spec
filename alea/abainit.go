package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponABAInit(signedABAInit *SignedMessage, abaInitMsgContainer *MsgContainer) error {
	
	fmt.Println("uponABAInit function")

	// get data
	abaInitData, err := signedABAInit.Message.GetABAInitData()
	if err != nil{
		errors.Wrap(err, "could not get abainitdata from signedABAInit")
	}

	// add the message to the container
	i.State.ABAState.ABAInitContainer.AddMsg(signedABAInit)
	abaInitMsgContainer.AddMsg(signedABAInit)


	// Increase counter and
	// if weak support (partial quorum) is achieved for a value different from the input, broadcast a new init
	// if strong support (quorum) is achieved, add to values and broadcast AUX
	if abaInitData.Vote == 1 {
		i.State.ABAState.Init1Counter += 1
	} else {
		i.State.ABAState.Init0Counter += 1
	}

	// weak support

	// if never sent INIT(1) and reached PartialQuorum (i.e. f+1, weak support), send INIT 1
	if !i.State.ABAState.SentInit1 && i.State.ABAState.Init1Counter >= i.State.Share.PartialQuorum{
		// send INIT(1)
		initMsg, err := CreateABAInit(i.State, i.config, byte(1), abaInitData.Round)
		if err != nil {
			errors.Wrap(err,"failed to create ABA Init message after weak support")
		}
		i.Broadcast(initMsg)
		// update sent flag
		i.State.ABAState.SentInit1 = true
	}
	// (same thing for 0) if never sent INIT(0) and reached PartialQuorum (i.e. f+1, weak support), send INIT 0
	if !i.State.ABAState.SentInit0 && i.State.ABAState.Init0Counter >= i.State.Share.PartialQuorum{
		// send INIT(0)
		initMsg, err := CreateABAInit(i.State, i.config, byte(0), abaInitData.Round)
		if err != nil {
			errors.Wrap(err,"failed to create ABA Init message after weak support")
		}
		i.Broadcast(initMsg)
		// update sent flag
		i.State.ABAState.SentInit0 = true
	}


	// strong support

	
	// if never sent AUX(1) and reached Quorum (i.e. 2f+1, strong support), sends AUX(1) and add 1 to values
	if  !i.State.ABAState.SentAux1 && i.State.ABAState.Init1Counter >= i.State.Share.Quorum {
		// sends AUX(1)
		auxMsg, err := CreateABAAux(i.State, i.config, byte(1), abaInitData.Round)
		if err != nil {
			errors.Wrap(err,"failed to create ABA Aux message after strong init support")
		}
		i.Broadcast(auxMsg)

		// initializes queue if it doesn't exist
		if _, exists := i.State.ABAState.Values[i.State.ABAState.Round]; !exists {
			i.State.ABAState.Values[i.State.ABAState.Round] = make([]byte,0)
		}
		// append vote
		i.State.ABAState.Values[i.State.ABAState.Round] = append(i.State.ABAState.Values[i.State.ABAState.Round],byte(1))

		// update sent flag
		i.State.ABAState.SentAux1 = true
	}
	// if never sent AUX(0) and reached Quorum (i.e. 2f+1, strong support), sends AUX(0) and add 0 to values
	if  !i.State.ABAState.SentAux0 && i.State.ABAState.Init0Counter >= i.State.Share.Quorum {
		// sends AUX(1)
		auxMsg, err := CreateABAAux(i.State, i.config, byte(0), abaInitData.Round)
		if err != nil {
			errors.Wrap(err,"failed to create ABA Aux message after strong init support")
		}
		i.Broadcast(auxMsg)

		// initializes queue if it doesn't exist
		if _, exists := i.State.ABAState.Values[i.State.ABAState.Round]; !exists {
			i.State.ABAState.Values[i.State.ABAState.Round] = make([]byte,0)
		}
		// append vote
		i.State.ABAState.Values[i.State.ABAState.Round] = append(i.State.ABAState.Values[i.State.ABAState.Round],byte(0))

		// update sent flag
		i.State.ABAState.SentAux0 = true
	}

	return nil

}
	

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
