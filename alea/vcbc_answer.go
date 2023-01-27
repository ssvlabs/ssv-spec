package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponVCBCAnswer(signedMessage *SignedMessage) error {

	// get Data
	vcbcAnswerData, err := signedMessage.Message.GetVCBCAnswerData()
	if err != nil {
		return errors.Wrap(err, "uponVCBCAnswer: could not get data from signedMessage")
	}

	// check if has local aggregated signature
	hasLocalSignature := i.State.VCBCState.hasU(vcbcAnswerData.Author, vcbcAnswerData.Priority)
	if hasLocalSignature {
		return nil
	}

	// update local values
	i.State.VCBCState.setU(vcbcAnswerData.Author, vcbcAnswerData.Priority, vcbcAnswerData.Proof)
	i.State.VCBCState.setM(vcbcAnswerData.Author, vcbcAnswerData.Priority, vcbcAnswerData.Proposals)

	// add vcbc output
	i.AddVCBCOutput(vcbcAnswerData.Proposals, vcbcAnswerData.Priority, vcbcAnswerData.Author)

	return nil
}

func CreateVCBCAnswer(state *State, config IConfig, proposals []*ProposalData, priority Priority, proof types.Signature, author types.OperatorID) (*SignedMessage, error) {
	vcbcAnswerData := &VCBCAnswerData{
		Proposals: proposals,
		Priority:  priority,
		Proof:     proof,
		Author:    author,
	}
	dataByts, err := vcbcAnswerData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCAnswer: could not encode vcbcAnswerData")
	}
	msg := &Message{
		MsgType:    VCBCAnswerMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCAnswer: failed signing filler msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
