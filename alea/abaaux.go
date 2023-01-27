package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponABAAux(signedABAAux *SignedMessage, abaAuxMsgContainer *MsgContainer) error {
	if i.verbose {
		fmt.Println("uponABAAux")
	}
	// get message Data
	ABAAuxData, err := signedABAAux.Message.GetABAAuxData()
	if err != nil {
		errors.Wrap(err, "uponABAAux: could not get ABAAuxData from signedABAAux")
	}

	// add the message to the containers
	abaAuxMsgContainer.AddMsg(signedABAAux)

	voteInLocalValues := false
	for _, value := range i.State.ABAState.Values[ABAAuxData.Round] {
		if value == ABAAuxData.Vote {
			voteInLocalValues = true
		}
	}
	if i.verbose {
		fmt.Println("\tvote received is in local values:", voteInLocalValues, ". Local values (of round", ABAAuxData.Round, "):", i.State.ABAState.Values[ABAAuxData.Round], ". Vote:", ABAAuxData.Vote)
	}

	if voteInLocalValues {
		// increment counter
		i.State.ABAState.AuxCounter[ABAAuxData.Vote] += 1
		if i.verbose {
			fmt.Println("\tincremented aux counter. Vote:", ABAAuxData.Vote, ". Conter:", i.State.ABAState.AuxCounter)
		}
	}

	// if received 2f+1 AUX messages, try to send CONF
	if (i.State.ABAState.AuxCounter[0]+i.State.ABAState.AuxCounter[1]) >= i.State.Share.Quorum && !i.State.ABAState.SentConf {
		if i.verbose {
			fmt.Println("\tgot quorum of AUX and never sent conf")
		}

		// broadcast CONF message
		confMsg, err := CreateABAConf(i.State, i.config, i.State.ABAState.Values[ABAAuxData.Round], ABAAuxData.Round)
		if err != nil {
			errors.Wrap(err, "uponABAAux: failed to create ABA Conf message after strong support")
		}
		if i.verbose {
			fmt.Println("\tbroadcasting ABAConf")
		}
		i.Broadcast(confMsg)

		// update sent flag
		i.State.ABAState.SentConf = true
	}

	return nil
}

func CreateABAAux(state *State, config IConfig, vote byte, round Round) (*SignedMessage, error) {
	ABAAuxData := &ABAAuxData{
		Vote:  vote,
		Round: round,
	}
	dataByts, err := ABAAuxData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAAux: could not encode abaaux data")
	}
	msg := &Message{
		MsgType:    ABAAuxMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAAux: failed signing abaaux msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
