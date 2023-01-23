package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func CreateABA(state *State, config IConfig, vote byte, round Round) (*SignedMessage, error) {
	abaData := &ABAData{
		Vote:			vote,
		Round:			round,					
	}
	dataByts, err := abaData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode aba data")
	}
	msg := &Message{
		MsgType:    ABAMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing aba msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
