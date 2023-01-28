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

	// if message is old, return
	if ABAConfData.Round < i.State.ABAState.Round {
		return nil
	}

	// sender
	senderID := signedABAConf.GetSigners()[0]

	alreadyReceived := i.State.ABAState.hasConf(ABAConfData.Round, senderID)
	if i.verbose {
		fmt.Println("\tsenderID:", senderID, ", votes:", ABAConfData.Votes, ", round:", ABAConfData.Round, ", already received before:", alreadyReceived)
	}
	// if already received this msg, return
	if alreadyReceived {
		return nil
	}

	// determine if votes list is contained in local round values list
	isContained := i.State.ABAState.isContainedInValues(ABAConfData.Round, ABAConfData.Votes)
	// list is contained -> update CONF counter
	if isContained {
		i.State.ABAState.setConf(ABAConfData.Round, senderID)
		if i.verbose {
			fmt.Println("\tupdated confcounter:", i.State.ABAState.countConf(ABAConfData.Round))
		}
	}

	// reached strong support -> try to decide value
	if i.State.ABAState.countConf(ABAConfData.Round) >= i.State.Share.Quorum {
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

			i.State.ABAState.setVInput(ABAConfData.Round+1, s)
			if i.verbose {
				fmt.Println("\tlength of values is 2", i.State.ABAState.Values[ABAConfData.Round], "-> storing coin to next Vin")
			}
		} else {
			if i.verbose {
				fmt.Println("\tlength of values is 1:", i.State.ABAState.Values[ABAConfData.Round])
			}
			i.State.ABAState.setVInput(ABAConfData.Round+1, i.State.ABAState.GetValues(ABAConfData.Round)[0])

			// if value has only one value, sends FINISH
			if i.State.ABAState.GetValues(ABAConfData.Round)[0] == s {
				if i.verbose {
					fmt.Println("\tvalue equal to S")
				}
				// check if indeed never sent FINISH message for this vote
				if !i.State.ABAState.sentFinish(s) {
					finishMsg, err := CreateABAFinish(i.State, i.config, s)
					if err != nil {
						errors.Wrap(err, "uponABAConf: failed to create ABA Finish message")
					}
					if i.verbose {
						fmt.Println("\tSending ABAFinish")
					}
					i.Broadcast(finishMsg)
					i.State.ABAState.setSentFinish(s, true)
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
		initMsg, err := CreateABAInit(i.State, i.config, i.State.ABAState.getVInput(i.State.ABAState.Round), i.State.ABAState.Round)
		if err != nil {
			errors.Wrap(err, "uponABAConf: failed to create ABA Init message")
		}
		if i.verbose {
			fmt.Println("\tSending ABAInit with new Vin:", i.State.ABAState.Vin[i.State.ABAState.Round], ", for round:", i.State.ABAState.Round)
		}
		i.Broadcast(initMsg)
	}

	return nil
}

func isValidABAConf(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != ABAConfMsgType {
		return errors.New("msg type is not ABAConfMsgType")
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

	ABAConfData, err := signedMsg.Message.GetABAConfData()
	if err != nil {
		return errors.Wrap(err, "could not get ABAConfData data")
	}
	if err := ABAConfData.Validate(); err != nil {
		return errors.Wrap(err, "ABAConfData invalid")
	}

	// vote
	votes := ABAConfData.Votes
	for _, vote := range votes {
		if vote != 0 && vote != 1 {
			return errors.Wrap(err, "vote different than 0 and 1")
		}
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
