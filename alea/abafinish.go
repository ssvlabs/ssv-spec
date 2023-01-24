package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponABAFinish(signedABAFinish *SignedMessage, abaFinishMsgContainer *MsgContainer) error {
	
	fmt.Println("uponABAFinish function")

	ABAFinishData, err := signedABAFinish.Message.GetABAFinishData()
	if err != nil{
		errors.Wrap(err, "could not get ABAFinishData from signedABAConf")
	}

	// add the message to the container
	i.State.ABAState.ABAFinishContainer.AddMsg(signedABAFinish)

	vote := ABAFinishData.Vote
	
	if vote == 1 {
		i.State.ABAState.Finish1Counter += 1
	} else {
		i.State.ABAState.Finish0Counter += 1
	}

	if !i.State.ABAState.SentFinish1 && i.State.ABAState.Finish1Counter >= i.State.Share.PartialQuorum {
		finishMsg, err := CreateABAFinish(i.State, i.config, byte(1))
		if err != nil {
			errors.Wrap(err,"failed to create ABA Finish message")
		}
		i.Broadcast(finishMsg)
		i.State.ABAState.SentFinish1 = true
	}
	if !i.State.ABAState.SentFinish0 && i.State.ABAState.Finish0Counter >= i.State.Share.PartialQuorum {
		finishMsg, err := CreateABAFinish(i.State, i.config, byte(0))
		if err != nil {
			errors.Wrap(err,"failed to create ABA Finish message")
		}
		i.Broadcast(finishMsg)
		i.State.ABAState.SentFinish0 = true
	}

	if i.State.ABAState.Finish1Counter >= i.State.Share.Quorum {
		i.State.ABAState.Vdecided = byte(1)
		i.State.ABAState.Terminate <- true
	}
	if i.State.ABAState.Finish0Counter >= i.State.Share.Quorum {
		i.State.ABAState.Vdecided = byte(0)
		i.State.ABAState.Terminate <- true
	}

	return nil

}


func CreateABAFinish(state *State, config IConfig, vote byte) (*SignedMessage, error) {
	ABAFinishData := &ABAFinishData{
		Vote:		vote,			
	}
	dataByts, err := ABAFinishData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode abafinish data")
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
		return nil, errors.Wrap(err, "failed signing abafinish msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
