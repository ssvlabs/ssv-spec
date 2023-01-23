package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func CreateVCBC(state *State, config IConfig,proposalData []*ProposalData, priority int) (*SignedMessage, error) {
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
		Round:      state.Round,
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
