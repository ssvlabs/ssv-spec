package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// LateCommitPastInstance tests process commit msg for a previously decided instance
func LateCommitPastInstance() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	msgs := testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], 1, testingutils.Testing4SharesSet())

	return &tests.ControllerSpecTest{
		Name: "late commit past instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				InputMessages:      msgs[0:7],
				DecidedVal:         []byte{1, 2, 3, 4},
				DecidedCnt:         1,
				ControllerPostRoot: "4904e750939440bf885052e33dadc77369fe4a942cbe9940bf4ec6c52baac1b7",
			},
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: append(msgs[7:14], testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], 4, &qbft.Message{
					MsgType:    qbft.CommitMsgType,
					Height:     qbft.FirstHeight,
					Round:      qbft.FirstRound,
					Identifier: identifier[:],
					Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
				})),
				DecidedVal:         []byte{1, 2, 3, 4},
				DecidedCnt:         1,
				ControllerPostRoot: "ab417f75ea3def610f96ae88a1406c757b973fb273cac6a6c434398b56d06283",
			},
		},
	}
}
