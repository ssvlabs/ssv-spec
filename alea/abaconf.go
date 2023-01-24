package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponABAConf(signedABAConf *SignedMessage, abaConfMsgContainer *MsgContainer) error {
	
	fmt.Println("uponABAConf function")

	ABAConfData, err := signedABAConf.Message.GetABAConfData()
	if err != nil{
		errors.Wrap(err, "could not get ABAConfData from signedABAConf")
	}

	// add the message to the container
	i.State.ABAState.ABAConfContainer.AddMsg(signedABAConf)

	votes := ABAConfData.Votes
	equalValues := 0
	for _, vote := range votes {
		for _, value := range i.State.ABAState.Values[ABAConfData.Round] {
			if vote == value {
				equalValues += 1
			}
		}
	}
	if equalValues == len(votes) {
		i.State.ABAState.ConfCounter += 1
	}

	if i.State.ABAState.ConfCounter >= i.State.Share.Quorum {
		s := i.State.ABAState.Coin(i.State.ABAState.Round)
		if len(i.State.ABAState.Values[ABAConfData.Round]) == 2 {
			i.State.ABAState.Vin = s
		} else {
			i.State.ABAState.Vin = i.State.ABAState.Values[ABAConfData.Round][0]
			if i.State.ABAState.Values[ABAConfData.Round][0] == s {
				finishMsg, err := CreateABAFinish(i.State, i.config, s)
				if err != nil {
					errors.Wrap(err,"failed to create ABA Finish message")
				}
				i.Broadcast(finishMsg)
				
				if s == 1 {
					i.State.ABAState.SentFinish1 = true
				} else {
					i.State.ABAState.SentFinish0 = true
				}
			}
		}

		i.State.ABAState.IncrementRound()

		initMsg, err := CreateABAInit(i.State, i.config, i.State.ABAState.Vin, i.State.ABAState.Round)
		if err != nil {
			errors.Wrap(err,"failed to create ABA Init message")
		}
		i.Broadcast(initMsg)
	}

	return nil

}



func CreateABAConf(state *State, config IConfig, votes []byte, round Round) (*SignedMessage, error) {
	ABAConfData := &ABAConfData{
		Votes:		votes,
		Round:		round,					
	}
	dataByts, err := ABAConfData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode abaconf data")
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
		return nil, errors.Wrap(err, "failed signing abaconf msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
