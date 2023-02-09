package hbbft

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
		return errors.Wrap(err, "uponABAAux: could not get ABAAuxData from signedABAAux")
	}

	i.State.HBBFTState.GetRound()

	// old message -> ignore
	if ABAAuxData.ACSRound < i.State.HBBFTState.GetRound() {
		return nil
	}
	if ABAAuxData.Round < i.State.HBBFTState.GetCurrentABAState().Round {
		return nil
	}

	// if future round -> intialize future state
	if ABAAuxData.ACSRound > i.State.HBBFTState.GetRound() {
		i.State.HBBFTState.InitMaps(ABAAuxData.ACSRound)
	}
	if ABAAuxData.Round > i.State.HBBFTState.GetABAState(ABAAuxData.ACSRound).Round {
		i.State.HBBFTState.GetABAState(ABAAuxData.ACSRound).InitializeRound(ABAAuxData.Round)
	}

	abaState := i.State.HBBFTState.GetABAState(ABAAuxData.ACSRound)

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
		abaState.setAux(ABAAuxData.Round, senderID, ABAAuxData.Vote)
		if i.verbose {
			fmt.Println("\tincremented aux counter. Vote:", ABAAuxData.Vote)
		}
	}

	// if received 2f+1 AUX messages, try to send CONF
	if (abaState.countAux(ABAAuxData.Round, 0)+abaState.countAux(ABAAuxData.Round, 1)) >= i.State.Share.Quorum && !abaState.sentConf(ABAAuxData.Round) {
		if i.verbose {
			fmt.Println("\tgot quorum of AUX and never sent conf")
		}
		if i.verbose {
			fmt.Println("\tcalculating q")
		}
		q := abaState.CountAuxInValues(ABAAuxData.Round)
		if i.verbose {
			fmt.Println("\tq:", q, ", quorum:", i.State.Share.Quorum)
		}

		if q < i.State.Share.Quorum {
			if i.verbose {
				fmt.Println("\tcurrent quorum of msgs doesn't reach quorum of msgs with votes in local values. q (votes in local values):", q)
			}
			return nil
		}

		if i.verbose {
			fmt.Println("\tsending:", abaState.Values[ABAAuxData.Round], "for round:", ABAAuxData.Round)
		}
		// broadcast CONF message
		confMsg, err := CreateABAConf(i.State, i.config, abaState.Values[ABAAuxData.Round], ABAAuxData.Round, ABAAuxData.ACSRound)
		if err != nil {
			return errors.Wrap(err, "uponABAAux: failed to create ABA Conf message after strong support")
		}
		if i.verbose {
			fmt.Println("\tbroadcasting ABAConf")
		}
		i.Broadcast(confMsg)

		// update sent flag
		abaState.setSentConf(ABAAuxData.Round, true)
		// process own conf msg
		i.uponABAConf(confMsg)
	}

	return nil
}

func isValidABAAux(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
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
		return errors.New("vote different than 0 and 1")
	}

	return nil
}

func CreateABAAux(state *State, config IConfig, vote byte, round Round, acRound ACSRound) (*SignedMessage, error) {
	ABAAuxData := &ABAAuxData{
		Vote:     vote,
		Round:    round,
		ACSRound: acRound,
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
