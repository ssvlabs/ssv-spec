package hbbft

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponRBCEcho(signedMsg *SignedMessage) error {

	return nil
}

func isValidRBCEcho(
	state *State,
	config IConfig,
	signedProposal *SignedMessage,
	operators []*types.Operator,
) error {

	return nil
}

// CreateRBCEcho
func CreateRBCEcho(state *State, config IConfig, markleTree []byte, branch []byte, erasureShare []byte) (*SignedMessage, error) {
	rbcEchoData := &RBCEchoData{
		MarkleTree:   markleTree,
		Branch:       branch,
		ErasureShare: erasureShare,
	}
	dataByts, err := rbcEchoData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode proposal data")
	}
	msg := &Message{
		MsgType:    RBCEchoMsgType,
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
