package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostFinish  tests a valid SignedPartialSignatureMessage post finished runner
func PostFinish() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	errCode := types.NoRunningDutyErrorCode
	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus post finish",
		testdoc.PostConsensusPostFinishDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name: "proposer",
				Runner: finishRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
			},
			{
				Name: "proposer (blinded block)",
				Runner: finishRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
			},
		},
		ks,
	)

	// Aggregator committee duty
	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name: "sync committee contribution",
		Runner: finishAggregatorCommitteeRunner(
			testingutils.AggregatorCommitteeRunner(ks),
			testingutils.TestingSyncCommitteeContributionDuty,
			testingutils.TestSyncCommitteeContributionConsensusData,
		),
		Duty: testingutils.TestingSyncCommitteeContributionDuty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[4], 4, ks))),
		},
		DontStartDuty:     true,
		ExpectedErrorCode: errCode,
	})
	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name: fmt.Sprintf("aggregator (%s)", version.String()),
				Runner: finishAggregatorCommitteeRunner(
					testingutils.AggregatorCommitteeRunner(ks),
					testingutils.TestingAggregatorDuty(version),
					testingutils.TestAggregatorConsensusData(version),
				),
				Duty: testingutils.TestingAggregatorDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[4], 4, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
			},
			{
				Name: fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Runner: finishAggregatorCommitteeRunner(
					testingutils.AggregatorCommitteeRunner(ks),
					testingutils.TestingAggregatorCommitteeDutyMixed(version),
					testingutils.TestAggregatorCommitteeConsensusData(version),
				),
				Duty: testingutils.TestingAggregatorCommitteeDutyMixed(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMixedMsg(ks.Shares[4], 4, version, ks))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
			},
		}...)
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{

			{
				Name: fmt.Sprintf("attester (%s)", version.String()),
				Runner: finishCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[4], 4, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
			},
			{
				Name: fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: finishCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingSyncCommitteeDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingSyncCommitteeDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[4], 4, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
			},
			{
				Name: fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: finishCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterAndSyncCommitteeDuties(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[4], 4, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
			},
		}...)
	}

	return multiSpecTest
}
