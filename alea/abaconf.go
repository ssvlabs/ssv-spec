package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponABAConf(signedABAConf *SignedMessage, abaConfMsgContainer *MsgContainer) error {
	
	fmt.Println("uponABAConf function")

	// get data
	ABAConfData, err := signedABAConf.Message.GetABAConfData()
	if err != nil{
		errors.Wrap(err, "could not get ABAConfData from signedABAConf")
	}

	// add the message to the container
	i.State.ABAState.ABAConfContainer.AddMsg(signedABAConf)
	abaConfMsgContainer.AddMsg(signedABAConf)

	// get votes in CONF message
	votes := ABAConfData.Votes

	// determine if votes list is contained in local round values list

	// determine the number of equal values
	equalValues := 0
	for _, vote := range votes {
		for _, value := range i.State.ABAState.Values[ABAConfData.Round] {
			if vote == value {
				equalValues += 1
			}
		}
	}
	// if number of equal values == length of list -> list is contained -> update CONF counter 
	if equalValues == len(votes) {
		i.State.ABAState.ConfCounter += 1
	}

	// reached strong support -> try to decide value
	if i.State.ABAState.ConfCounter >= i.State.Share.Quorum {

		// get common coin
		s := i.State.ABAState.Coin(i.State.ABAState.Round)

		// if values = {0,1}, choose randomly (i.e. coin) value for next round
		if len(i.State.ABAState.Values[ABAConfData.Round]) == 2 {
			i.State.ABAState.Vin = s
		} else {
			i.State.ABAState.Vin = i.State.ABAState.Values[ABAConfData.Round][0]

			// if value has only one value, sends FINISH
			if i.State.ABAState.Values[ABAConfData.Round][0] == s {
				// check if indeed never sent FINISH message for this vote
				if ((s == 1 && !i.State.ABAState.SentFinish1) || (s == 0 && !i.State.ABAState.SentFinish0)) {
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
		}

		// increment round
		i.State.ABAState.IncrementRound()

		// start new round sending INIT message with vote
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
