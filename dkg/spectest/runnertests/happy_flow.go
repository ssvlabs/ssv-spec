package runnertests

import (
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/dkg/types"
	dkgtypes "github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	index := 4
	ks := testingutils.Testing4SharesSet()
	dataSet := testutils.TestDepositSignDataSetFourOperators()
	pre := testutils.TestRunner(dataSet, *ks, uint64(index))
	return &MsgProcessingSpecTest{
		Name:   "happy flow",
		Pre:    pre,
		KeySet: ks,
		Output: dataSet.MakeOutput(dkgtypes.OperatorID(index)),
		Messages: []*types.Message{
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, testutils.PlaceholderMessage()).(*types.Message),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, dataSet.ParsedPartialSigMessage(2)).(*types.Message),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, dataSet.ParsedPartialSigMessage(3)).(*types.Message),
		},
	}
}
