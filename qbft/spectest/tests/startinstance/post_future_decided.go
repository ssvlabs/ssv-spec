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
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	return &tests.ControllerSpecTest{
		Name: "start instance prev decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         inputData,
				InputMessages:      testingutils.DecidingMsgsForHeight(inputData, identifier, qbft.FirstHeight, testingutils.Testing4SharesSet()),
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "df9b2787df60e1e15b0c840410592d27803d44cd5fb086cfa8fc23181cea6293",
			},
			{
				InputValue:         inputData,
				ControllerPostRoot: "4fe3611d5d04b8b34ab9dd470e352dd9693c9bd3b340a8b42d0bc608fb4d59bc",
			},
		},
	}
}
