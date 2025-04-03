package consensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidDecided13Operators tests a valid decided value (13 operators)
func ValidDecided13Operators() tests.SpecTest {
	ks := testingutils.Testing13SharesSet()
	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid decided 13 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee contribution",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.RoleSyncCommitteeContribution),
				PostDutyRunnerStateRoot: "2feec534f4ef14ea96e83b17d4feaf3f40e4739a72fc7118c7da1ef2c1331937",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb), ks, types.RoleProposer),
				PostDutyRunnerStateRoot: "942d09a2aaaddf2ef7b2388fd1a9ffb94d590c645fc90c2c6b27fa5414afd61e",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
			},
			{
				Name:                    "proposer (blinded block)",
				Runner:                  testingutils.ProposerBlindedBlockRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb), ks, types.RoleProposer),
				PostDutyRunnerStateRoot: "abf4eaf532eda91fd6db48136938985a655174ab5101e698d71737c0ed48c8c3",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
			},
		},
	}

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:                    fmt.Sprintf("aggregator (%s)", version.String()),
			Runner:                  testingutils.AggregatorRunner(ks),
			Duty:                    testingutils.TestingAggregatorDuty(version),
			Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData(version), ks, types.RoleAggregator),
			PostDutyRunnerStateRoot: "dd79a9da2d025a61d4e4ebdebbc9a331a876a43bb8c8e7c274132ae1c4a35175",
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version),
				testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1, version),
			},
		},
		)
	}

	for _, version := range testingutils.SupportedAttestationVersions {

		height := qbft.Height(testingutils.TestingDutySlotV(version))

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:     fmt.Sprintf("attester (%s)", version.String()),
				Runner:   testingutils.CommitteeRunner(ks),
				Duty:     testingutils.TestingAttesterDuty(version),
				Messages: testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, height),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, version),
				},
			},
			{
				Name:     fmt.Sprintf("sync committee (%s)", version.String()),
				Runner:   testingutils.CommitteeRunner(ks),
				Duty:     testingutils.TestingSyncCommitteeDuty(version),
				Messages: testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, height),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1, version),
				},
			},
			{
				Name:     fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner:   testingutils.CommitteeRunner(ks),
				Duty:     testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, height),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version),
				},
			},
		}...)
	}

	return multiSpecTest
}
