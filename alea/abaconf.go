package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponABAConf(signedABAConf *SignedMessage, abaConfMsgContainer *MsgContainer) error {
	if i.verbose {
		fmt.Println("uponABAConf")
	}
	// get data
	ABAConfData, err := signedABAConf.Message.GetABAConfData()
	if err != nil {
		errors.Wrap(err, "uponABAConf:could not get ABAConfData from signedABAConf")
	}

	// add the message to the container
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
	if i.verbose {
		fmt.Println("\tnum equal values:", equalValues, ". Len of votes:", len(votes))
	}
	// if number of equal values == length of list -> list is contained -> update CONF counter
	if equalValues == len(votes) {
		i.State.ABAState.ConfCounter += 1
		if i.verbose {
			fmt.Println("\tupdated confcounter:", i.State.ABAState.ConfCounter)
		}
	}

	// reached strong support -> try to decide value
	if i.State.ABAState.ConfCounter >= i.State.Share.Quorum {
		if i.verbose {
			fmt.Println("\treached quorum")
		}

		// get common coin
		s := i.State.ABAState.Coin(i.State.ABAState.Round)
		if i.verbose {
			fmt.Println("\tcoin:", s)
		}

		// if values = {0,1}, choose randomly (i.e. coin) value for next round
		if len(i.State.ABAState.Values[ABAConfData.Round]) == 2 {
			i.State.ABAState.Vin = s
			if i.verbose {
				fmt.Println("\tlength of values is 2", i.State.ABAState.Values[ABAConfData.Round], "-> storing coin to next Vin")
			}
		} else {
			if i.verbose {
				fmt.Println("\tlength of values is 1:", i.State.ABAState.Values[ABAConfData.Round])
			}
			i.State.ABAState.Vin = i.State.ABAState.Values[ABAConfData.Round][0]

			// if value has only one value, sends FINISH
			if i.State.ABAState.Values[ABAConfData.Round][0] == s {
				if i.verbose {
					fmt.Println("\tvalue equal to S")
				}
				// check if indeed never sent FINISH message for this vote
				if !i.State.ABAState.SentFinish[s] {
					finishMsg, err := CreateABAFinish(i.State, i.config, s)
					if err != nil {
						errors.Wrap(err, "uponABAConf: failed to create ABA Finish message")
					}
					if i.verbose {
						fmt.Println("\tSending ABAFinish")
					}
					i.Broadcast(finishMsg)
					i.State.ABAState.SentFinish[s] = true
					if i.verbose {
						fmt.Println("\tupdated SentFinish:", i.State.ABAState.SentFinish)
					}
				}
			}
		}

		// increment round
		if i.verbose {
			fmt.Println("\twill icrement round. Round now:", i.State.ABAState.Round)
		}
		i.State.ABAState.IncrementRound()
		if i.verbose {
			fmt.Println("\tnew round:", i.State.ABAState.Round)
		}

		// start new round sending INIT message with vote
		initMsg, err := CreateABAInit(i.State, i.config, i.State.ABAState.Vin, i.State.ABAState.Round)
		if err != nil {
			errors.Wrap(err, "uponABAConf: failed to create ABA Init message")
		}
		if i.verbose {
			fmt.Println("\tSending ABAInit with new Vin:", i.State.ABAState.Vin, ", for round:", i.State.ABAState.Round)
		}
		i.Broadcast(initMsg)
	}

	return nil
}

func CreateABAConf(state *State, config IConfig, votes []byte, round Round) (*SignedMessage, error) {
	ABAConfData := &ABAConfData{
		Votes: votes,
		Round: round,
	}
	dataByts, err := ABAConfData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAConf: could not encode abaconf data")
	}
	msg := &Message{
		MsgType:    ABAConfMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAConf: failed signing abaconf msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
