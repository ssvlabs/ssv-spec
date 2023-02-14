package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponVCBCSend(signedMessage *SignedMessage) error {

	if i.verbose {
		fmt.Println("uponVCBCSend")
	}

	// get Data
	vcbcSendData, err := signedMessage.Message.GetVCBCSendData()
	if err != nil {
		errors.New("uponVCBCSend: could not get vcbcSendData data from signedMessage")
	}

	// check if it was already received. If not -> store
	if !i.State.VCBCState.HasM(vcbcSendData.Author, vcbcSendData.Priority) {
		i.State.VCBCState.setM(vcbcSendData.Author, vcbcSendData.Priority, vcbcSendData.Proposals)
	}

	// get sender of the message
	senderID := signedMessage.GetSigners()[0]
	if i.verbose {
		fmt.Println("\tgot senderID:", senderID)
	}

	// if it's the sender and the author, just add the ready signature
	if senderID == i.State.Share.OperatorID && senderID == vcbcSendData.Author {
		i.AddOwnVCBCReady(vcbcSendData.Proposals, vcbcSendData.Priority)
		return nil
	}

	// if the Author of the VCBC is the same as the sender of the message -> sign and answer with READY
	if senderID == vcbcSendData.Author {
		if i.verbose {
			fmt.Println("\tsenderID is the same as the author")
		}

		hash, err := GetProposalsHash(vcbcSendData.Proposals)
		if err != nil {
			return errors.New("uponVCBCSend: could not get hash of proposals")
		}
		if i.verbose {
			fmt.Println("\tgot hash")
		}

		// create VCBCReady message with proof
		vcbcReadyMsg, err := CreateVCBCReady(i.State, i.config, hash, vcbcSendData.Priority, vcbcSendData.Author)
		if err != nil {
			return errors.New("uponVCBCSend: failed to create VCBCReady message with proof")
		}

		if i.verbose {
			fmt.Println("\tBroadcasting VCBC ready")
		}
		// FIX ME : send specifically to author
		i.Broadcast(vcbcReadyMsg)
	}

	return nil
}

func isValidVCBCSend(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != VCBCSendMsgType {
		return errors.New("msg type is not VCBCSend")
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

	VCBCSendData, err := signedMsg.Message.GetVCBCSendData()
	if err != nil {
		return errors.Wrap(err, "could not get vcbcsend data")
	}
	if err := VCBCSendData.Validate(); err != nil {
		return errors.Wrap(err, "VCBCSendData invalid")
	}

	// author
	author := VCBCSendData.Author
	authorInCommittee := false
	for _, opID := range operators {
		if opID.OperatorID == author {
			authorInCommittee = true
		}
	}
	if !authorInCommittee {
		return errors.New("author (OperatorID) doesn't exist in Committee")
	}

	if author != signedMsg.GetSigners()[0] {
		return errors.New("author of VCBCSend differs from sender of the message")
	}

	// priority
	priority := VCBCSendData.Priority
	if state.VCBCState.HasM(author, priority) {
		if !state.VCBCState.EqualM(author, priority, VCBCSendData.Proposals) {
			return errors.New("existing (priority,author) with different proposals")
		}
	}

	return nil
}

func CreateVCBCSend(state *State, config IConfig, proposals []*ProposalData, priority Priority, author types.OperatorID) (*SignedMessage, error) {
	vcbcSendData := &VCBCSendData{
		Proposals: proposals,
		Priority:  priority,
		Author:    author,
	}
	dataByts, err := vcbcSendData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCSend: could not encode vcbcSendData")
	}
	msg := &Message{
		MsgType:    VCBCSendMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCSend: failed signing filler msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
