package hbbft

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponABAFinish(signedABAFinish *SignedMessage) error {
	if i.verbose {
		fmt.Println("uponABAFinish")
	}
	// get data
	ABAFinishData, err := signedABAFinish.Message.GetABAFinishData()
	if err != nil {
		return errors.Wrap(err, "uponABAFinish: could not get ABAFinishData from signedABAConf")
	}

	// old message -> ignore
	if ABAFinishData.ACSRound < i.State.HBBFTState.GetRound() {
		return nil
	}
	// if future round -> intialize future state
	if ABAFinishData.ACSRound > i.State.HBBFTState.GetRound() {
		i.State.HBBFTState.InitMaps(ABAFinishData.ACSRound)
	}

	abaState := i.State.HBBFTState.GetABAState(ABAFinishData.ACSRound)

	// add the message to the container
	abaState.ABAFinishContainer.AddMsg(signedABAFinish)

	// sender
	senderID := signedABAFinish.GetSigners()[0]

	alreadyReceived := abaState.hasFinish(senderID)
	if i.verbose {
		fmt.Println("\tsenderID:", senderID, ", vote:", ABAFinishData.Vote, ", already received before:", alreadyReceived)
	}
	// if never received this msg, update
	if !alreadyReceived {

		// get vote from FINISH message
		vote := ABAFinishData.Vote

		// increment counter
		abaState.setFinish(senderID, vote)
		if i.verbose {
			fmt.Println("\tincremented finish counter:", abaState.FinishCounter)
		}
		if i.verbose {
			fmt.Println("\tSentFinish:", abaState.SentFinish)
		}
	}

	// if FINISH(b) reached partial quorum and never broadcasted FINISH(b), broadcast
	if !abaState.sentFinish(byte(0)) && !abaState.sentFinish(byte(1)) {
		for _, vote := range []byte{0, 1} {

			if abaState.countFinish(vote) >= i.State.Share.PartialQuorum {
				if i.verbose {
					fmt.Println("\treached partial quorum of finish and never sent -> sending new, for vote:", vote)
					fmt.Println("\tsentFinish[vote]:", abaState.sentFinish(vote), ", vote", vote)

				}
				// broadcast FINISH
				finishMsg, err := CreateABAFinish(i.State, i.config, vote, ABAFinishData.ACSRound)
				if err != nil {
					return errors.Wrap(err, "uponABAFinish: failed to create ABA Finish message")
				}
				if i.verbose {
					fmt.Println("\tsending ABAFinish")
				}
				i.Broadcast(finishMsg)

				// update sent flag
				abaState.setSentFinish(vote, true)
				// process own finish msg
				i.uponABAFinish(finishMsg)
			}
		}
	}

	// if FINISH(b) reached Quorum, decide for b and send termination signal
	for _, vote := range []byte{0, 1} {
		if abaState.countFinish(vote) >= i.State.Share.Quorum {
			if i.verbose {
				fmt.Println("\treached quorum for vote:", vote)
			}
			abaState.setDecided(vote)
			abaState.setTerminate(true)
		}
	}

	return nil
}

func isValidABAFinish(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
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
		return errors.New("vote different than 0 and 1")
	}

	return nil
}

func CreateABAFinish(state *State, config IConfig, vote byte, acsRound ACSRound) (*SignedMessage, error) {
	ABAFinishData := &ABAFinishData{
		Vote:     vote,
		ACSRound: acsRound,
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
