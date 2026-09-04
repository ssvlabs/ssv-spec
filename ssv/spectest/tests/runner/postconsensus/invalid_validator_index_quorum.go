package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidValidatorIndexQuorum tests a quorum of post consensus messages with an invalid validator index
func InvalidValidatorIndexQuorum() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	expectedError := types.UnknownValidatorIndexErrorCode
	sccSlot := testingutils.TestingSyncCommitteeContributionDuty.Slot
	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus invalid validator index quorum",
		testdoc.PostConsensusInvalidValidatorIndexQuorumDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name: "proposer",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongValidatorIndexProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongValidatorIndexProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongValidatorIndexProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedError,
			},
			{
				Name: "proposer (blinded block)",
				Runner: decideRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongValidatorIndexProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongValidatorIndexProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongValidatorIndexProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedError,
			},
		},
		ks,
	)

	// Aggregator committee duty
	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name: "sync committee contribution",
		Runner: decideAggregatorCommitteeRunner(
			testingutils.AggregatorCommitteeRunner(ks),
			testingutils.TestingSyncCommitteeContributionDuty,
			testingutils.TestSyncCommitteeContributionConsensusData,
		),
		Duty: testingutils.TestingSyncCommitteeContributionDuty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusWrongValidatorIndexSyncCommitteeContributionMsg(ks.Shares[1], 1, ks, sccSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusWrongValidatorIndexSyncCommitteeContributionMsg(ks.Shares[2], 2, ks, sccSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusWrongValidatorIndexSyncCommitteeContributionMsg(ks.Shares[3], 3, ks, sccSlot))),
		},
		DontStartDuty: true,
		// No error is expected as AggregatorCommitteeRunner doesn't validate the validator indexes
	})
	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name: fmt.Sprintf("aggregator (%s)", version.String()),
				Runner: decideAggregatorCommitteeRunner(
					testingutils.AggregatorCommitteeRunner(ks),
					testingutils.TestingAggregatorDuty(version),
					testingutils.TestAggregatorConsensusData(version),
				),
				Duty: testingutils.TestingAggregatorDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusWrongValidatorIndexAggregatorMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusWrongValidatorIndexAggregatorMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusWrongValidatorIndexAggregatorMsg(ks.Shares[3], 3, version))),
				},
				DontStartDuty: true,
				// No error is expected as AggregatorCommitteeRunner doesn't validate the validator indexes
			},
			{
				Name: fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Runner: decideAggregatorCommitteeRunner(
					testingutils.AggregatorCommitteeRunner(ks),
					testingutils.TestingAggregatorCommitteeDutyMixed(version),
					testingutils.TestAggregatorCommitteeConsensusData(version),
				),
				Duty: testingutils.TestingAggregatorCommitteeDutyMixed(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMixedMsgWrongValIdx(ks.Shares[1], 1, version, ks))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMixedMsgWrongValIdx(ks.Shares[2], 2, version, ks))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMixedMsgWrongValIdx(ks.Shares[3], 3, version, ks))),
				},
				DontStartDuty: true,
				// No error is expected as AggregatorCommitteeRunner doesn't validate the validator indexes
			},
		}...)
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{

			{
				Name: fmt.Sprintf("attester (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongValidatorIndexAttestationMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongValidatorIndexAttestationMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongValidatorIndexAttestationMsg(ks.Shares[3], 3, version))),
				},
				DontStartDuty: true,
				// No error is expected for the CommitteeRunner since we don't assume that operators must be synced on the validators set
			},
			{
				Name: fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingSyncCommitteeDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingSyncCommitteeDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongValidatorIndexSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongValidatorIndexSyncCommitteeMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongValidatorIndexSyncCommitteeMsg(ks.Shares[3], 3, version))),
				},
				DontStartDuty: true,
				// No error is expected for the CommitteeRunner since we don't assume that operators must be synced on the validators set
			},
			{
				Name: fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterAndSyncCommitteeDuties(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongValidatorIndexAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongValidatorIndexAttestationAndSyncCommitteeMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongValidatorIndexAttestationAndSyncCommitteeMsg(ks.Shares[3], 3, version))),
				},
				DontStartDuty: true,
				// No error is expected for the CommitteeRunner since we don't assume that operators must be synced on the validators set
			},
		}...)
	}

	return multiSpecTest
}
