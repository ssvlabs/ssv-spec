package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponVCBCRequest(signedMessage *SignedMessage) error {

	// get Data
	vcbcRequestData, err := signedMessage.Message.GetVCBCRequestData()
	if err != nil {
		return errors.Wrap(err, "uponVCBCRequest: could not get data from signedMessage")
	}

	// check if has local aggregated signature. If not, return
	if !i.State.VCBCState.hasU(vcbcRequestData.Author, vcbcRequestData.Priority) {
		return nil
	}

	proposals := i.State.VCBCState.GetM(vcbcRequestData.Author, vcbcRequestData.Priority)
	u := i.State.VCBCState.GetU(vcbcRequestData.Author, vcbcRequestData.Priority)

	msgToBroadcast, err := CreateVCBCAnswer(i.State, i.config, proposals, vcbcRequestData.Priority, u, vcbcRequestData.Author)
	if err != nil {
		return errors.Wrap(err, "uponVCBCRequest: failed to create VCBCAnswer message")
	}

	// FIX ME : send only to requester
	i.Broadcast(msgToBroadcast)

	return nil
}

func isValidVCBCRequest(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != VCBCRequestMsgType {
		return errors.New("msg type is not VCBCRequestMsgType")
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

	VCBCRequestData, err := signedMsg.Message.GetVCBCRequestData()
	if err != nil {
		return errors.Wrap(err, "could not get VCBCRequestData data")
	}
	if err := VCBCRequestData.Validate(); err != nil {
		return errors.Wrap(err, "VCBCRequestData invalid")
	}

	// author
	author := VCBCRequestData.Author
	authorInCommittee := false
	for _, opID := range operators {
		if opID.OperatorID == author {
			authorInCommittee = true
		}
	}
	if !authorInCommittee {
		return errors.New("author (OperatorID) doesn't exist in Committee")
	}

	return nil
}

func CreateVCBCRequest(state *State, config IConfig, priority Priority, author types.OperatorID) (*SignedMessage, error) {
	vcbcRequestData := &VCBCRequestData{
		Priority: priority,
		Author:   author,
	}
	dataByts, err := vcbcRequestData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCRequest: could not encode vcbcRequestData")
	}
	msg := &Message{
		MsgType:    VCBCRequestMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCRequest: failed signing filler msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
