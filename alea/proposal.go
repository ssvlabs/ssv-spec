package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)

// uponProposal process proposal message
// Assumes proposal message is valid!
func (i *Instance) uponProposal(signedProposal *SignedMessage, proposeMsgContainer *MsgContainer) error {
    
	fmt.Println("uponProposal function")

	// get Data
	proposalDataReceived, err := signedProposal.Message.GetProposalData()
	if err != nil{
		errors.Wrap(err, "could not get proposal data from signedProposal")
	}

	// check if message has been already delivered
	if i.State.S.hasProposal(proposalDataReceived) {
		fmt.Println("\tmessage already contained in S queue (already delivered)")
		return nil
	}

	// Add message to container
    proposeMsgContainer.AddMsg(signedProposal)
	fmt.Println("\tAdded message to container")

    // Check if container has less maximum size. If so, returns
    if proposeMsgContainer.Len(i.State.AleaDefaultRound) < i.config.GetBatchSize() {
		fmt.Println("\tReturning. Len not big enough to match batch size")
		fmt.Println("\tcontainer length:",proposeMsgContainer.Len(i.State.AleaDefaultRound),", batch size:",i.config.GetBatchSize())
		return nil
    }

	fmt.Println("\treached batch size")

	// broadcast VCBC with received proposals

	// get messages
	signedMessages := proposeMsgContainer.MessagesForRound(i.State.AleaDefaultRound)

	// extract proposal datas
	proposalData := make([]*ProposalData,0)
	for i := range signedMessages {
		data, err := signedMessages[i].Message.GetProposalData();
		if err != nil {
			errors.Wrap(err, "could not get proposal data from message in container")
		}
		proposalData = append(proposalData,data)
	}
	fmt.Println("\tcreated proposal data")

	// create VCBC message and broadcasts
	msgToBroadcast, err := CreateVCBC(i.State, i.config, proposalData, i.State.Priority)
	fmt.Println("\tcreated VCBC message to broadcast")
	if err != nil {
		return errors.Wrap(err, "failed to create VCBC message")
	}
    if err := i.Broadcast(msgToBroadcast); err != nil {
        return errors.Wrap(err, "failed to broadcast VCBC message")
    }
	fmt.Println("\tbroadcasted")

    // Clear container
    proposeMsgContainer.Clear()
	fmt.Println("\tcleared container")

	// Increment priority
	i.State.Priority += 1
	fmt.Println("\tincremented priority")
    return nil
}


func isValidProposal(
	state *State,
	config IConfig,
	signedProposal *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedProposal.Message.MsgType != ProposalMsgType {
		return errors.New("msg type is not proposal")
	}
	if signedProposal.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedProposal.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if err := signedProposal.Signature.VerifyByOperators(signedProposal, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}
	if !signedProposal.MatchedSigners([]types.OperatorID{proposer(state, config, signedProposal.Message.Round)}) {
		return errors.New("proposal leader invalid")
	}

	proposalData, err := signedProposal.Message.GetProposalData()
	if err != nil {
		return errors.Wrap(err, "could not get proposal data")
	}
	if err := proposalData.Validate(); err != nil {
		return errors.Wrap(err, "proposalData invalid")
	}

	// if err := isProposalJustification(
	// 	state,
	// 	config,
	// 	proposalData.RoundChangeJustification,
	// 	proposalData.PrepareJustification,
	// 	state.Height,
	// 	signedProposal.Message.Round,
	// 	proposalData.Data,
	// 	valCheck,
	// ); err != nil {
	// 	return errors.Wrap(err, "proposal not justified")
	// }

	if (state.ProposalAcceptedForCurrentRound == nil && signedProposal.Message.Round == state.Round) ||
		signedProposal.Message.Round > state.Round {
		return nil
	}
	return errors.New("proposal is not valid with current state")
}

// isProposalJustification returns nil if the proposal and round change messages are valid and justify a proposal message for the provided round, value and leader
func isProposalJustification(
	state *State,
	config IConfig,
	roundChangeMsgs []*SignedMessage,
	prepareMsgs []*SignedMessage,
	height Height,
	round Round,
	value []byte,
	valCheck ProposedValueCheckF,
) error {
	if err := valCheck(value); err != nil {
		return errors.Wrap(err, "proposal value invalid")
	}

	if round == FirstRound {
		return nil
	} else {
		// check all round changes are valid for height and round
		// no quorum, duplicate signers,  invalid still has quorum, invalid no quorum
		// prepared
		// for _, rc := range roundChangeMsgs {
		// 	if err := validRoundChange(state, config, rc, height, round); err != nil {
		// 		return errors.Wrap(err, "change round msg not valid")
		// 	}
		// }

		// check there is a quorum
		if !HasQuorum(state.Share, roundChangeMsgs) {
			return errors.New("change round has no quorum")
		}

		// previouslyPreparedF returns true if any on the round change messages have a prepared round and value
		previouslyPrepared, err := func(rcMsgs []*SignedMessage) (bool, error) {
			// for _, rc := range rcMsgs {
			// 	rcData, err := rc.Message.GetRoundChangeData()
			// 	if err != nil {
			// 		return false, errors.Wrap(err, "could not get round change data")
			// 	}
			// 	if rcData.Prepared() {
			// 		return true, nil
			// 	}
			// }
			return false, nil
		}(roundChangeMsgs)
		if err != nil {
			return errors.Wrap(err, "could not calculate if previously prepared")
		}

		if !previouslyPrepared {
			return nil
		} else {

			// check prepare quorum
			if !HasQuorum(state.Share, prepareMsgs) {
				return errors.New("prepares has no quorum")
			}

			// get a round change data for which there is a justification for the highest previously prepared round
			// rcm, err := highestPrepared(roundChangeMsgs)
			// if err != nil {
			// 	return errors.Wrap(err, "could not get highest prepared")
			// }
			// if rcm == nil {
			// 	return errors.New("no highest prepared")
			// }
			// rcmData, err := rcm.Message.GetRoundChangeData()
			// if err != nil {
			// 	return errors.Wrap(err, "could not get round change data")
			// }

			// proposed value must equal highest prepared value
			// if !bytes.Equal(value, rcmData.PreparedValue) {
			// 	return errors.New("proposed data doesn't match highest prepared")
			// }

			// validate each prepare message against the highest previously prepared value and round
			// for _, pm := range prepareMsgs {
			// 	if err := validSignedPrepareForHeightRoundAndValue(
			// 		config,
			// 		pm,
			// 		height,
			// 		rcmData.PreparedRound,
			// 		rcmData.PreparedValue,
			// 		state.Share.Committee,
			// 	); err != nil {
			// 		return errors.New("signed prepare not valid")
			// 	}
			// }
			return nil
		}
	}
}

func proposer(state *State, config IConfig, round Round) types.OperatorID {
	// TODO - https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/29ae5a44551466453a84d4d17b9e083ecf189d97/dafny/spec/L1/node_auxiliary_functions.dfy#L304-L323
	return config.GetProposerF()(state, round)
}

// CreateProposal
/**
  	Proposal(
                        signProposal(
                            UnsignedProposal(
                                |current.blockchain|,
                                newRound,
                                digest(block)),
                            current.id),
                        block,
                        extractSignedRoundChanges(roundChanges),
                        extractSignedPrepares(prepares));
*/
func CreateProposal(state *State, config IConfig, value []byte/*, roundChanges, prepares []*SignedMessage*/) (*SignedMessage, error) {
	proposalData := &ProposalData{
		Data:                     value,
		// RoundChangeJustification: roundChanges,
		// PrepareJustification:     prepares,
	}
	dataByts, err := proposalData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode proposal data")
	}
	msg := &Message{
		MsgType:    ProposalMsgType,
		Height:     state.Height,
		Round:      state.AleaDefaultRound,
		Identifier: state.ID,
		Data:       dataByts,
	}
	
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing prepare msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
