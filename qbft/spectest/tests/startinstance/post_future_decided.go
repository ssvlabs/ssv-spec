package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreviousDecided tests starting an instance when the previous one decided
func PreviousDecided() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "start instance prev decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], qbft.FirstHeight, testingutils.Testing4SharesSet()),
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: []byte{1, 2, 3, 4},
					DecidedCnt: 1,
				},
				ControllerPostRoot: "0370be5066cbbf1efead61d9b182309afd989b3b720163f7029cbad79537eb4b",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "0e0db8b36601bef53328c1f1e94e4c3faa02084c743cce138dc2edcce2e5d79e",
			},
		},
	}
}
