package validation

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/types/testingutils"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// NonDecidedWithMultipleSigners tests a non-decided message with multiple signers
func NonDecidedWithMultipleSigners() tests.SpecTest {

	msgF := func(msgType qbft.MessageType) *types.SignedSSVMessage {
		return &types.SignedSSVMessage{
			OperatorIDs: []types.OperatorID{1, 2, 3},
			Signatures:  [][]byte{testingutils.MockRSASignature[:], testingutils.MockRSASignature[:], testingutils.MockRSASignature[:]},
			SSVMessage: &types.SSVMessage{
				MsgType: types.SSVConsensusMsgType,
				MsgID:   testingutils.DefaultMsgID,
				Data: testingutils.EncodeQbftMessage(&qbft.Message{
					MsgType: msgType,
					Round:   qbft.FirstRound,
				}),
			},
		}
	}

	possibleMessageTypes := []qbft.MessageType{qbft.ProposalMsgType, qbft.PrepareMsgType, qbft.RoundChangeMsgType}

	// Create multi test
	multiTest := &MultiMessageValidationTest{
		Name:  "non decided with multiple signers",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Add test cases
	for _, msgType := range possibleMessageTypes {
		multiTest.Tests = append(multiTest.Tests, &MessageValidationTest{
			Name:          fmt.Sprintf("%v type", msgType),
			Messages:      [][]byte{testingutils.EncodeMessage(msgF(msgType))},
			ExpectedError: validation.ErrNonDecidedWithMultipleSigners.Error(),
		})
	}

	return multiTest
}
