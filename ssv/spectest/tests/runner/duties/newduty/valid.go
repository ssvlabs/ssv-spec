package newduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests a valid start duty
func Valid() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &MultiStartNewRunnerDutySpecTest{
		Name: "new duty valid",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "29862cc6054edc8547efcb5ae753290971d664b9c39768503b4d66e1b52ecb06",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "56eafcb33392ded888a0fefe30ba49e52aa00ab36841cb10c9dc1aa2935af347",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
		},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*StartNewRunnerDutySpecTest{
			{
				Name:                    fmt.Sprintf("attester (%s)", version.String()),
				Runner:                  testingutils.CommitteeRunner(ks),
				Duty:                    testingutils.TestingAttesterDuty(version),
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "7f926e54651ed34901256e8c82a40658647afe17cb089f6c1d7406e7350f4c2e",
				OutputMessages:          []*types.PartialSignatureMessages{},
			},
			{
				Name:                    fmt.Sprintf("sync committee (%s)", version.String()),
				Runner:                  testingutils.CommitteeRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeDuty(version),
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "29862cc6054edc8547efcb5ae753290971d664b9c39768503b4d66e1b52ecb06",
				OutputMessages:          []*types.PartialSignatureMessages{},
			},
			{
				Name:                    fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner:                  testingutils.CommitteeRunner(ks),
				Duty:                    testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "29862cc6054edc8547efcb5ae753290971d664b9c39768503b4d66e1b52ecb06",
				OutputMessages:          []*types.PartialSignatureMessages{},
			},
		}...)
	}

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &StartNewRunnerDutySpecTest{
			Name:      fmt.Sprintf("aggregator (%s)", version.String()),
			Runner:    testingutils.AggregatorRunner(ks),
			Duty:      testingutils.TestingAggregatorDuty(version),
			Threshold: ks.Threshold,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version), // broadcasts when starting a new duty
			},
		})
	}
	return multiSpecTest

}
