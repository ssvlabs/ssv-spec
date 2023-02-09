package hbbft

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponDEC(signedMsg *SignedMessage) error {

	return nil
}

func isValidDEC(
	state *State,
	config IConfig,
	signedProposal *SignedMessage,
	operators []*types.Operator,
) error {

	return nil
}

// CreateDEC
func CreateDEC(state *State, config IConfig, acsRound ACSRound, operatorID types.OperatorID, sender types.OperatorID, decryptedShare []byte) (*SignedMessage, error) {
	decData := &DECData{
		ACSRound:       acsRound,
		OperatorID:     operatorID,
		Sender:         sender,
		DecryptedShare: decryptedShare,
	}
	dataByts, err := decData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode proposal data")
	}
	msg := &Message{
		MsgType:    DECMsgType,
		Height:     state.Height,
		Round:      HBBFTDefaultRound,
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
