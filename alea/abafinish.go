package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponABAFinish(signedABAFinish *SignedMessage, abaFinishMsgContainer *MsgContainer) error {
	if i.verbose {
		fmt.Println("uponABAFinish")
	}
	// get data
	ABAFinishData, err := signedABAFinish.Message.GetABAFinishData()
	if err != nil {
		errors.Wrap(err, "uponABAFinish: could not get ABAFinishData from signedABAConf")
	}

	// add the message to the container
	abaFinishMsgContainer.AddMsg(signedABAFinish)

	// sender
	senderID := signedABAFinish.GetSigners()[0]

	alreadyReceived := i.State.ABAState.hasFinish(senderID)
	if i.verbose {
		fmt.Println("\tsenderID:", senderID, ", vote:", ABAFinishData.Vote, ", already received before:", alreadyReceived)
	}
	// if already received this msg, return
	if alreadyReceived {
		return nil
	}

	// get vote from FINISH message
	vote := ABAFinishData.Vote

	// increment counter
	i.State.ABAState.setFinish(senderID, vote)
	if i.verbose {
		fmt.Println("\tincremented finish counter:", i.State.ABAState.FinishCounter)
	}

	if i.verbose {
		fmt.Println("\tSentFinish:", i.State.ABAState.SentFinish)
	}

	// if FINISH(b) reached partial quorum and never broadcasted FINISH(b), broadcast
	for _, vote := range []byte{0, 1} {

		if !i.State.ABAState.sentFinish(vote) && i.State.ABAState.countFinish(vote) >= i.State.Share.PartialQuorum {
			if i.verbose {
				fmt.Println("\treached partial quorum of finish and never sent -> sending new, for vote:", vote)
				fmt.Println("\tsentFinish[vote]:", i.State.ABAState.sentFinish(vote), ", vote", vote)

			}
			// broadcast FINISH
			finishMsg, err := CreateABAFinish(i.State, i.config, vote)
			if err != nil {
				errors.Wrap(err, "uponABAFinish: failed to create ABA Finish message")
			}
			if i.verbose {
				fmt.Println("\tsending ABAFinish")
			}
			i.Broadcast(finishMsg)
			// update sent flag
			i.State.ABAState.setSentFinish(vote, true)
		}
	}

	// if FINISH(b) reached Quorum, decide for b and send termination signal
	for _, vote := range []byte{0, 1} {
		if i.State.ABAState.countFinish(vote) >= i.State.Share.Quorum {
			if i.verbose {
				fmt.Println("\treached quorum for vote:", vote)
			}
			i.State.ABAState.Vdecided = vote
			i.State.ABAState.Terminate = true
		}
	}

	return nil
}

func isValidABAFinish(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != ABAFinishMsgType {
		return errors.New("msg type is not ABAFinishMsgType")
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

	ABAFinishData, err := signedMsg.Message.GetABAFinishData()
	if err != nil {
		return errors.Wrap(err, "could not get ABAFinishData data")
	}
	if err := ABAFinishData.Validate(); err != nil {
		return errors.Wrap(err, "ABAFinishData invalid")
	}

	// vote
	vote := ABAFinishData.Vote
	if vote != 0 && vote != 1 {
		return errors.Wrap(err, "vote different than 0 and 1")
	}

	return nil
}

func CreateABAFinish(state *State, config IConfig, vote byte) (*SignedMessage, error) {
	ABAFinishData := &ABAFinishData{
		Vote: vote,
	}
	dataByts, err := ABAFinishData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode abafinish data")
	}
	msg := &Message{
		MsgType:    ABAFinishMsgType,
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
