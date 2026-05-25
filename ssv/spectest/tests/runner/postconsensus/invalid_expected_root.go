package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidExpectedRoot tests 1 expected root which doesn't match the signed root
func InvalidExpectedRoot() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	expectedErrorCode := types.WrongSigningRootErrorCode
	sccSlot := testingutils.TestingSyncCommitteeContributionDuty.Slot
	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus invalid expected roots",
		testdoc.PostConsensusInvalidExpectedRootDoc,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
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
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusWrongSyncCommitteeContributionMsg(ks.Shares[1], 1, ks, sccSlot))),
		},
		DontStartDuty: true,
		// No error is expected as AggregatorCommitteeRunner doesn't validate the roots
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusWrongAggregatorMsg(ks.Shares[1], 1, version))),
				},
				DontStartDuty: true,
				// No error is expected as AggregatorCommitteeRunner doesn't validate the roots
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMixedWrongMsg(ks.Shares[1], 1, version, ks))),
				},
				DontStartDuty: true,
				// No error is expected as AggregatorCommitteeRunner doesn't validate the roots
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongAttestationMsg(ks.Shares[1], 1, version))),
				},
				DontStartDuty: true,
				// No error for this committee duty
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongSyncCommitteeMsg(ks.Shares[1], 1, version))),
				},
				DontStartDuty: true,
				// No error for this committee duty
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version))),
				},
				DontStartDuty: true,
				// No error for this committee duty
			},
		}...)
	}

	return multiSpecTest
}
