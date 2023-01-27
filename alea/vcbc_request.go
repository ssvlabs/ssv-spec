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

	proposals := i.State.VCBCState.getM(vcbcRequestData.Author, vcbcRequestData.Priority)
	u := i.State.VCBCState.getU(vcbcRequestData.Author, vcbcRequestData.Priority)

	msgToBroadcast, err := CreateVCBCAnswer(i.State, i.config, proposals, vcbcRequestData.Priority, u, vcbcRequestData.Author)
	if err != nil {
		return errors.Wrap(err, "uponVCBCRequest: failed to create VCBCAnswer message")
	}

	// FIX ME : send only to requester
	i.Broadcast(msgToBroadcast)

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
