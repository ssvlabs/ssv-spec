package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreviousDecided tests starting an instance when the previous one decided
func PreviousDecided() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
	return &tests.ControllerSpecTest{
		Name: "start instance prev decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         inputData,
				InputMessages:      testingutils.DecidingMsgsForHeight(inputData, identifier, qbft.FirstHeight, testingutils.Testing4SharesSet()),
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "e7823a17225ee7f1163e71b0fc0b67df888cfe287f5ec7a6454ab105a402a998",
			},
			{
				InputValue:         inputData,
				ControllerPostRoot: "c484a6a6539781a02ff749fbf5d50a06000d069e151bcd56d6d81359f7217e2c",
			},
		},
	}
}
