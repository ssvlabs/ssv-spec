package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
	"fmt"
)


func (i *Instance) uponVCBC(signedVCBC *SignedMessage, vcbcMsgContainer *MsgContainer) error {
    
	fmt.Println("uponVCBC function")

	// OperatorID who sent the message (and first signed it)
	senderID := signedVCBC.GetSigners()[0]

	fmt.Println("\tgot senderID:",senderID)

	// initializes queue if it doesn't exist
	if _, exists := i.State.queues[senderID]; !exists {
		i.State.queues[senderID] = NewVCBCQueue()
	}

	// gets the sender's associated queue
	queue := i.State.queues[senderID]
	fmt.Println("\tgot senderID's queue:",queue)


	// gets the VCBC data from the message
	vcbcData, err := signedVCBC.Message.GetVCBCData()
	if err != nil {
    	return errors.Wrap(err, "could not get VCBCData from message")
	}
	fmt.Println("\tgot vcbcData:",vcbcData)


	// check if it was already delivered
	if i.State.S.hasProposalList(vcbcData.ProposalData) {
		fmt.Println("\tlist of proposals from VCBC already contained in S queue (already delivered)")
		return nil
	}
	fmt.Println("\tlocal S hasn't yet this list of proposal:",vcbcData)


	queue.Enqueue(vcbcData.ProposalData, vcbcData.Priority)
	fmt.Println("\tenqueueing proposal list and priority")
	fmt.Println("\tnew queue:",i.State.queues[senderID])

	return nil
}


func isValidVCBC(
	state *State,
	config IConfig,
	signedVCBC *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedVCBC.Message.MsgType != VCBCMsgType {
		return errors.New("msg type is not proposal")
	}
	if signedVCBC.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedVCBC.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if err := signedVCBC.Signature.VerifyByOperators(signedVCBC, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	vcbcData, err := signedVCBC.Message.GetVCBCData()
	if err != nil {
		return errors.Wrap(err, "could not get vcbc data")
	}
	if err := vcbcData.Validate(); err != nil {
		return errors.Wrap(err, "vcbcData invalid")
	}

	// if err := isProposalJustification(
	// 	state,
	// 	config,
	// 	proposalData.RoundChangeJustification,
	// 	proposalData.PrepareJustification,
	// 	state.Height,
	// 	signedVCBC.Message.Round,
	// 	proposalData.Data,
	// 	valCheck,
	// ); err != nil {
	// 	return errors.Wrap(err, "proposal not justified")
	// }

	if (state.ProposalAcceptedForCurrentRound == nil && signedVCBC.Message.Round == state.Round) ||
		signedVCBC.Message.Round > state.Round {
		return nil
	}
	return errors.New("proposal is not valid with current state")
}


func CreateVCBC(state *State, config IConfig, proposalData []*ProposalData, priority Priority) (*SignedMessage, error) {
	vcbcData := &VCBCData{
		ProposalData:	proposalData,
		Priority:		priority,					
	}
	dataByts, err := vcbcData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode vcbc data")
	}
	msg := &Message{
		MsgType:    VCBCMsgType,
		Height:     state.Height,
		Round:      state.AleaDefaultRound,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing vcbc msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
