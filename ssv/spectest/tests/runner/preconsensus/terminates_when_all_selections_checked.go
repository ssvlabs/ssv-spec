package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TerminatesWhenAllSelectionsChecked ensures the runner terminates after testing the selections for all possible duties
func TerminatesWhenAllSelectionsChecked() tests.SpecTest {
	// Create key share map for several validators
	ks := testingutils.Testing4SharesSet()
	numValidators := 10
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus terminates when all selections checked",
		testdoc.PreConsensusAggregatorCommitteeAllSelectionsCheckedAndNoAggregatorsDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {

		// Duty
		mixedDuty := testingutils.TestingAggregatorCommitteeDuty(validatorsIndexList, []int{}, version)

		// Create beacon node that selects no aggregator
		aggregatorsNotToBeSelectedMap := make(map[phase0.CommitteeIndex]bool)
		aggregatorsNotToBeSelectedKeys := make([]phase0.CommitteeIndex, 0)
		aggregatorsNotToBeSelectedValues := make([]bool, 0)
		for _, duty := range mixedDuty.ValidatorDuties {
			if duty.Type == types.BNRoleAggregator {
				aggregatorsNotToBeSelectedMap[duty.CommitteeIndex] = false
				aggregatorsNotToBeSelectedKeys = append(aggregatorsNotToBeSelectedKeys, duty.CommitteeIndex)
				aggregatorsNotToBeSelectedValues = append(aggregatorsNotToBeSelectedValues, false)
			}
		}
		beacon := testingutils.NewTestingBeaconNode()
		beacon.SetAggregators(aggregatorsNotToBeSelectedMap)

		runner, err := testingutils.ConstructBaseRunnerWithShareMapAndBeaconNode(types.RoleAggregatorCommittee, shareMap, beacon)
		if err != nil {
			panic(err)
		}

		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("terminates when all selections checked (%s)", version.String()),
			Runner: runner,
			Duty:   mixedDuty,
			Messages: []*types.SignedSSVMessage{
				// Send quorum that covers every duty, thus testing all selections
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 2))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 3))),

				// Post-consensus message raises no running duty error
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1, version))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1),
			},
			BeaconAggregators:       aggregatorsNotToBeSelectedKeys,
			BeaconAggregatorsValues: aggregatorsNotToBeSelectedValues,
			// Post-consensus message raises no running duty error as runner should have terminated after testing all selections
			ExpectedErrorCode: types.NoRunningDutyErrorCode,
		})
	}

	return multiSpecTest
}
