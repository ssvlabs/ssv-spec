package randao

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/processmsg/_basemsgtests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func BaseTests() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	testsToRun := make([]*tests.MsgProcessingSpecTest, 0)
	for _, baseTest := range _basemsgtests.AllBaseTests {
		generatedTest := baseTest.BuildTest(types.BNRoleProposer)
		ssvMessages := make([]*types.SSVMessage, 0)
		for _, m := range generatedTest.Msgs {
			ssvMessages = append(ssvMessages, testingutils.SSVMsgProposer(nil, m))
		}

		testsToRun = append(testsToRun, &tests.MsgProcessingSpecTest{
			Name:                    generatedTest.Name,
			Runner:                  testingutils.ProposerRunner(ks),
			Duty:                    testingutils.TestingProposerDuty,
			Messages:                ssvMessages,
			PostDutyRunnerStateRoot: generatedTest.PostDutyRunnerStateRoot,
			OutputMessages: []*ssv.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
			},
			ExpectedError: generatedTest.ExpectedError,
		})
	}
	return &tests.MultiMsgProcessingSpecTest{
		Name:  "randao base tests",
		Tests: testsToRun,
	}
}
