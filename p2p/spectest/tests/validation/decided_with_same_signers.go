package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DecidedWithSameSigners tests sending two decided messages with the same signers
func DecidedWithSameSigners() tests.SpecTest {

	// Function to create a decided message
	decidedMessage := func(signers []types.OperatorID, fullData []byte) *types.SignedSSVMessage {

		root, err := qbft.HashDataRoot(fullData)
		if err != nil {
			panic(err)
		}
		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Round:      qbft.FirstRound,
				Root:       root,
				Identifier: testingutils.DefaultMsgID[:],
			}),
		}
		signatures := [][]byte{}
		for _, signer := range signers {
			signatures = append(signatures, testingutils.SignSSVMessage(testingutils.DefaultKeySet.OperatorKeys[signer], ssvMsg))
		}
		return &types.SignedSSVMessage{
			OperatorIDs: signers,
			Signatures:  signatures,
			SSVMessage:  ssvMsg,
			FullData:    fullData,
		}
	}

	return &MessageValidationTest{
		Name: "decided with same signers",
		Messages: [][]byte{
			testingutils.EncodeMessage(decidedMessage([]types.OperatorID{1, 2, 3}, []byte{1})),
			testingutils.EncodeMessage(decidedMessage([]types.OperatorID{1, 2, 3}, []byte{2})),
		},
		ExpectedError: validation.ErrDecidedWithSameSigners.Error(),
	}
}
