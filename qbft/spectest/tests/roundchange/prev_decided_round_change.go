package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PrevDecidedRoundChange tests process round change msg for a decided instance
func PrevDecidedRoundChange() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "roudn change prev decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				InputMessages:      testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], qbft.FirstHeight, testingutils.Testing4SharesSet()),
				DecidedVal:         []byte{1, 2, 3, 4},
				DecidedCnt:         1,
				ControllerPostRoot: "4904e750939440bf885052e33dadc77369fe4a942cbe9940bf4ec6c52baac1b7",
			},
		},
	}
}
