package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponABAInit(signedABAInit *SignedMessage, abaInitMsgContainer *MsgContainer) error {

	if i.verbose {
		fmt.Println("uponABAInit")
	}

	// get data
	abaInitData, err := signedABAInit.Message.GetABAInitData()
	if err != nil {
		errors.Wrap(err, "uponABAInit: could not get abainitdata from signedABAInit")
	}

	// add the message to the container
	abaInitMsgContainer.AddMsg(signedABAInit)

	// Increase counter
	i.State.ABAState.InitCounter[abaInitData.Vote] += 1
	if i.verbose {
		fmt.Println("\tupdated counter. Vote:", abaInitData.Vote, ". InitCounter:", i.State.ABAState.InitCounter)
	}

	// weak support -> send INIT
	// if never sent INIT(b) but reached PartialQuorum (i.e. f+1, weak support), send INIT(b)
	for _, vote := range []byte{0, 1} {
		if !i.State.ABAState.SentInit[vote] && i.State.ABAState.InitCounter[vote] >= i.State.Share.PartialQuorum {
			if i.verbose {
				fmt.Println("\tgot weak support for (and never sent):", vote)
			}
			// send INIT
			initMsg, err := CreateABAInit(i.State, i.config, vote, abaInitData.Round)
			if err != nil {
				errors.Wrap(err, "uponABAInit: failed to create ABA Init message after weak support")
			}
			if i.verbose {
				fmt.Println("\tsending INIT")
			}
			i.Broadcast(initMsg)
			// update sent flag
			i.State.ABAState.SentInit[vote] = true
		}
	}

	// strong support -> send AUX
	// if never sent AUX(b) but reached Quorum (i.e. 2f+1, strong support), sends AUX(b) and add b to values
	for _, vote := range []byte{0, 1} {

		if !i.State.ABAState.SentAux[vote] && i.State.ABAState.InitCounter[vote] >= i.State.Share.Quorum {
			if i.verbose {
				fmt.Println("\tgot strong support and never sent AUX:", vote)
			}

			// initializes queue if it doesn't exist
			if _, exists := i.State.ABAState.Values[abaInitData.Round]; !exists {
				i.State.ABAState.Values[abaInitData.Round] = make([]byte, 0)
			}
			// append vote
			i.State.ABAState.Values[abaInitData.Round] = append(i.State.ABAState.Values[abaInitData.Round], vote)
			if i.verbose {
				fmt.Println("\tadded vote to local values for round", abaInitData.Round, ", values:", i.State.ABAState.Values[abaInitData.Round])
			}

			// sends AUX(b)
			auxMsg, err := CreateABAAux(i.State, i.config, vote, abaInitData.Round)
			if err != nil {
				errors.Wrap(err, "uponABAInit: failed to create ABA Aux message after strong init support")
			}
			if i.verbose {
				fmt.Println("\tsending ABAAux")
			}
			i.Broadcast(auxMsg)

			// update sent flag
			i.State.ABAState.SentAux[vote] = true
		}
	}

	return nil
}

func CreateABAInit(state *State, config IConfig, vote byte, round Round) (*SignedMessage, error) {
	ABAInitData := &ABAInitData{
		Vote:  vote,
		Round: round,
	}
	dataByts, err := ABAInitData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAInit: could not encode abainit data")
	}
	msg := &Message{
		MsgType:    ABAInitMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAInit: failed signing abainit msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
