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

	// if message is old, return
	if abaInitData.Round < i.State.ABAState.Round {
		return nil
	}

	// sender
	senderID := signedABAInit.GetSigners()[0]

	alreadyReceived := i.State.ABAState.hasInit(abaInitData.Round, senderID)
	if i.verbose {
		fmt.Println("\tsenderID:", senderID, ", vote:", abaInitData.Vote, ", round:", abaInitData.Round, ", already received before:", alreadyReceived)
	}
	// if already received this msg, return
	if alreadyReceived {
		return nil
	}

	// Set received msg
	i.State.ABAState.setInit(abaInitData.Round, senderID, abaInitData.Vote)

	if i.verbose {
		fmt.Println("\tupdated counter. Vote:", abaInitData.Vote, ". InitCounter:", i.State.ABAState.InitCounter)
	}

	// weak support -> send INIT
	// if never sent INIT(b) but reached PartialQuorum (i.e. f+1, weak support), send INIT(b)
	for _, vote := range []byte{0, 1} {
		if !i.State.ABAState.sentInit(abaInitData.Round, vote) && i.State.ABAState.countInit(abaInitData.Round, vote) >= i.State.Share.PartialQuorum {
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
			i.State.ABAState.setSentInit(abaInitData.Round, vote, true)
		}
	}

	// strong support -> send AUX
	// if never sent AUX(b) but reached Quorum (i.e. 2f+1, strong support), sends AUX(b) and add b to values
	for _, vote := range []byte{0, 1} {

		if !i.State.ABAState.sentAux(abaInitData.Round, vote) && i.State.ABAState.countInit(abaInitData.Round, vote) >= i.State.Share.Quorum {
			if i.verbose {
				fmt.Println("\tgot strong support and never sent AUX:", vote)
			}

			// append vote

			i.State.ABAState.AddToValues(abaInitData.Round, vote)
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
			i.State.ABAState.setSentAux(abaInitData.Round, vote, true)
		}
	}

	return nil
}

func isValidABAInit(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != ABAInitMsgType {
		return errors.New("msg type is not ABAInitMsgType")
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

	ABAInitData, err := signedMsg.Message.GetABAInitData()
	if err != nil {
		return errors.Wrap(err, "could not get ABAInitData data")
	}
	if err := ABAInitData.Validate(); err != nil {
		return errors.Wrap(err, "ABAInitData invalid")
	}

	// vote
	vote := ABAInitData.Vote
	if vote != 0 && vote != 1 {
		return errors.Wrap(err, "vote different than 0 and 1")
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
