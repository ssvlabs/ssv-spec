package hbbft

import (
	"encoding/json"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponTransaction(signedMsg *SignedMessage) error {

	return nil
}

func EncodeTransactions(transactions []*TransactionData) ([]byte, error) {
	return json.Marshal(transactions)
}

func isValidTransaction(
	state *State,
	config IConfig,
	signedProposal *SignedMessage,
	operators []*types.Operator,
) error {

	return nil
}

// CreateTransaction
func CreateTransaction(state *State, config IConfig, data []byte) (*SignedMessage, error) {
	transactionData := &TransactionData{
		Data: data,
	}
	dataByts, err := transactionData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode proposal data")
	}
	msg := &Message{
		MsgType:    TransactionMsgType,
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
