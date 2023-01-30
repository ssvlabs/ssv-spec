package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponABAAux(signedABAAux *SignedMessage) error {
	if i.verbose {
		fmt.Println("uponABAAux")
	}
	// get message Data
	ABAAuxData, err := signedABAAux.Message.GetABAAuxData()
	if err != nil {
		errors.Wrap(err, "uponABAAux: could not get ABAAuxData from signedABAAux")
	}

	// if future round -> intialize future state
	if ABAAuxData.ACRound > i.State.ACState.ACRound {
		i.State.ACState.InitializeRound(ABAAuxData.ACRound)
	}
	if ABAAuxData.Round > i.State.ACState.GetCurrentABAState().Round {
		i.State.ACState.GetCurrentABAState().InitializeRound(ABAAuxData.Round)
	}
	// old message -> ignore
	if ABAAuxData.ACRound < i.State.ACState.ACRound {
		return nil
	}
	if ABAAuxData.Round < i.State.ACState.GetCurrentABAState().Round {
		return nil
	}

	abaState := i.State.ACState.GetABAState(ABAAuxData.ACRound)

	// add the message to the containers
	abaState.ABAAuxContainer.AddMsg(signedABAAux)

	// sender
	senderID := signedABAAux.GetSigners()[0]

	alreadyReceived := abaState.hasAux(ABAAuxData.Round, senderID, ABAAuxData.Vote)
	if i.verbose {
		fmt.Println("\tsenderID:", senderID, ", vote:", ABAAuxData.Vote, ", round:", ABAAuxData.Round, ", already received before:", alreadyReceived)
	}
	// if never received this msg, increment counter
	if !alreadyReceived {
		voteInLocalValues := abaState.existsInValues(ABAAuxData.Round, ABAAuxData.Vote)
		if i.verbose {
			fmt.Println("\tvote received is in local values:", voteInLocalValues, ". Local values (of round", ABAAuxData.Round, "):", abaState.Values[ABAAuxData.Round], ". Vote:", ABAAuxData.Vote)
		}

		if voteInLocalValues {
			// increment counter

			abaState.setAux(ABAAuxData.Round, senderID, ABAAuxData.Vote)
			if i.verbose {
				fmt.Println("\tincremented aux counter. Vote:", ABAAuxData.Vote)
			}
		}
	}

	// if received 2f+1 AUX messages, try to send CONF
	if (abaState.countAux(ABAAuxData.Round, 0)+abaState.countAux(ABAAuxData.Round, 1)) >= i.State.Share.Quorum && !abaState.sentConf(ABAAuxData.Round) {
		if i.verbose {
			fmt.Println("\tgot quorum of AUX and never sent conf")
		}

		// broadcast CONF message
		confMsg, err := CreateABAConf(i.State, i.config, abaState.Values[ABAAuxData.Round], ABAAuxData.Round, ABAAuxData.ACRound)
		if err != nil {
			errors.Wrap(err, "uponABAAux: failed to create ABA Conf message after strong support")
		}
		if i.verbose {
			fmt.Println("\tbroadcasting ABAConf")
		}
		i.Broadcast(confMsg)

		// update sent flag
		abaState.setSentConf(ABAAuxData.Round, true)
		abaState.setConf(ABAAuxData.Round, i.State.Share.OperatorID)
	}

	return nil
}

func isValidABAAux(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != ABAAuxMsgType {
		return errors.New("msg type is not ABAAuxMsgType")
	}
	if signedMsg.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedMsg.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if err := signedMsg.Signature.VerifyByOperators(signedMsg, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	ABAAuxData, err := signedMsg.Message.GetABAAuxData()
	if err != nil {
		return errors.Wrap(err, "could not get ABAAuxData data")
	}
	if err := ABAAuxData.Validate(); err != nil {
		return errors.Wrap(err, "ABAAuxData invalid")
	}

	// vote
	vote := ABAAuxData.Vote
	if vote != 0 && vote != 1 {
		return errors.Wrap(err, "vote different than 0 and 1")
	}

	return nil
}

func CreateABAAux(state *State, config IConfig, vote byte, round Round, acRound ACRound) (*SignedMessage, error) {
	ABAAuxData := &ABAAuxData{
		Vote:    vote,
		Round:   round,
		ACRound: acRound,
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
