package validation

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicatedConsensusMessage tests sending a consensus message twice and triggering the message count error
func DuplicatedConsensusMessage() tests.SpecTest {

	// Function to get a consensus message for a given type
	msgForType := func(msgType qbft.MessageType) *types.SignedSSVMessage {
		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    msgType,
				Round:      qbft.FirstRound,
				Root:       testingutils.TestingQBFTRootData,
				Identifier: testingutils.DefaultMsgID[:],
			}),
		}

		fullData := []byte{}
		if msgType == qbft.ProposalMsgType {
			fullData = testingutils.TestingQBFTFullData
		}

		msg := &types.SignedSSVMessage{
			OperatorIDs: []types.OperatorID{1},
			Signatures:  [][]byte{testingutils.SignSSVMessage(testingutils.DefaultKeySet.OperatorKeys[1], ssvMsg)},
			SSVMessage:  ssvMsg,
			FullData:    fullData,
		}
		return msg
	}

	// Create multi test
	multiTest := &MultiMessageValidationTest{
		Name:  "duplicated consensus message",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Possible consensus message types
	messageTypes := []qbft.MessageType{qbft.ProposalMsgType, qbft.PrepareMsgType, qbft.CommitMsgType, qbft.RoundChangeMsgType}

	// Add test cases
	for _, msgType := range messageTypes {
		multiTest.Tests = append(multiTest.Tests, &MessageValidationTest{
			Name:          fmt.Sprintf("%v type", msgType),
			Messages:      [][]byte{testingutils.EncodeMessage(msgForType(msgType)), testingutils.EncodeMessage(msgForType(msgType))},
			ExpectedError: validation.ErrDuplicatedMessage.Error(),
		})
	}

	return multiTest
}
