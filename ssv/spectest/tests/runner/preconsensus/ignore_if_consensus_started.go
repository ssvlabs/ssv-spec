package preconsensus

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// IgnoreIfAlreadyStartedConsensus ensures the runner ignores a pre-consensus message if consensus has already started
func IgnoreIfAlreadyStartedConsensus() tests.SpecTest {

	// Create key share map for several validators
	ks := testingutils.Testing4SharesSet()
	numValidators := 10
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus ignores message if consensus already started",
		testdoc.PreConsensusAggregatorCommitteeIgnoreIfAlreadyStartedConsensusDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {

		// Duty with several validators
		mixedDuty := testingutils.TestingAggregatorCommitteeDuty(validatorsIndexList, []int{}, version)

		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("ignores message if consensus already started (%s)", version.String()),
			Runner: testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
			Duty:   mixedDuty,
			Messages: []*types.SignedSSVMessage{
				// Send quorum that covers every duty, thus testing all selections
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 2))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 3))),

				// This message should be ignored
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 4))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1),
			},
			ExpectedErrorCode: types.AggCommPreConsensusIgnoredSinceAlreadyStartedConsensusErrorCode,
		})
	}

	return multiSpecTest
}
