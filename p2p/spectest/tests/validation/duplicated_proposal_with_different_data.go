package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicatedProposalWithDifferentData tests sending two proposals messages for the same duty and round but with different full data values
func DuplicatedProposalWithDifferentData() tests.SpecTest {

	proposalMsg := func(fullData []byte) *types.SignedSSVMessage {
		root, err := qbft.HashDataRoot(fullData)
		if err != nil {
			panic(err)
		}
		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Round:      qbft.FirstRound,
				Root:       root,
				Identifier: testingutils.DefaultMsgID[:],
			}),
		}

		return &types.SignedSSVMessage{
			OperatorIDs: []types.OperatorID{1},
			Signatures:  [][]byte{testingutils.SignSSVMessage(testingutils.DefaultKeySet.OperatorKeys[1], ssvMsg)},
			SSVMessage:  ssvMsg,
			FullData:    fullData,
		}
	}

	return &MessageValidationTest{
		Name: "duplicated proposal with different data",
		Messages: [][]byte{
			testingutils.EncodeMessage(proposalMsg([]byte{1})),
			testingutils.EncodeMessage(proposalMsg([]byte{2})),
		},
		ExpectedError: validation.ErrDuplicatedProposalWithDifferentData.Error(),
	}
}
