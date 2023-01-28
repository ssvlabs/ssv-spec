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

	// if message is old, return
	if ABAAuxData.Round < i.State.ABAState.Round {
		return nil
	}

	// sender
	senderID := signedABAAux.GetSigners()[0]

	alreadyReceived := i.State.ABAState.hasAux(ABAAuxData.Round, senderID)
	if i.verbose {
		fmt.Println("\tsenderID:", senderID, ", vote:", ABAAuxData.Vote, ", round:", ABAAuxData.Round, ", already received before:", alreadyReceived)
	}
	// if already received this msg, return
	if alreadyReceived {
		return nil
	}

	voteInLocalValues := i.State.ABAState.existsInValues(ABAAuxData.Round, ABAAuxData.Vote)
	if i.verbose {
		fmt.Println("\tvote received is in local values:", voteInLocalValues, ". Local values (of round", ABAAuxData.Round, "):", i.State.ABAState.Values[ABAAuxData.Round], ". Vote:", ABAAuxData.Vote)
	}

	if voteInLocalValues {
		// increment counter

		i.State.ABAState.setAux(ABAAuxData.Round, senderID, ABAAuxData.Vote)
		if i.verbose {
			fmt.Println("\tincremented aux counter. Vote:", ABAAuxData.Vote)
		}
	}

	// if received 2f+1 AUX messages, try to send CONF
	if (i.State.ABAState.countAux(ABAAuxData.Round, 0)+i.State.ABAState.countAux(ABAAuxData.Round, 1)) >= i.State.Share.Quorum && !i.State.ABAState.sentConf(ABAAuxData.Round) {
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
		i.State.ABAState.setSentConf(ABAAuxData.Round, true)
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
