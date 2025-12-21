package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidBeaconSignatureInQuorum tests a quorum with an invalid PartialSignatureMessage signature
func InvalidBeaconSignatureInQuorum() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	expectedErrorCode := types.ReconstructSignatureErrorCode

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus invalid beacon signature in quorum",
		testdoc.PostConsensusInvalidBeaconSignatureInQuorumDoc,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongSigProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb))),
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongSigProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
			},
		},
		ks,
	)

	//// Aggregator committee duty
	//sccSlot := testingutils.TestingSyncCommitteeContributionDuty.Slot
	//multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
	//	Name: "sync committee contribution",
	//	Runner: decideAggregatorCommitteeRunner(
	//		testingutils.AggregatorCommitteeRunner(ks),
	//		testingutils.TestingSyncCommitteeContributionDuty,
	//		testingutils.TestSyncCommitteeContributionConsensusData,
	//	),
	//	Duty: testingutils.TestingSyncCommitteeContributionDuty,
	//	Messages: []*types.SignedSSVMessage{
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusWrongSigSyncCommitteeContributionMsg(ks.Shares[1], 1, ks, sccSlot))),
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks))),
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks))),
	//	},
	//	DontStartDuty:     true,
	//	ExpectedErrorCode: expectedErrorCode,
	//})
	//for _, version := range testingutils.SupportedAggregatorVersions {
	//	multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
	//		{
	//			Name: fmt.Sprintf("aggregator (%s)", version.String()),
	//			Runner: decideAggregatorCommitteeRunner(
	//				testingutils.AggregatorCommitteeRunner(ks),
	//				testingutils.TestingAggregatorDuty(version),
	//				testingutils.TestAggregatorConsensusData(version),
	//			),
	//			Duty: testingutils.TestingAggregatorDuty(version),
	//			Messages: []*types.SignedSSVMessage{
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusWrongSigAggregatorMsg(ks.Shares[1], 1, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3, version))),
	//			},
	//			DontStartDuty:     true,
	//			ExpectedErrorCode: expectedErrorCode,
	//		},
	//		{
	//			Name: fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
	//			Runner: decideAggregatorCommitteeRunner(
	//				testingutils.AggregatorCommitteeRunner(ks),
	//				testingutils.TestingAggregatorCommitteeDutyMixed(version),
	//				testingutils.TestAggregatorCommitteeConsensusData(version),
	//			),
	//			Duty: testingutils.TestingAggregatorCommitteeDutyMixed(version),
	//			Messages: []*types.SignedSSVMessage{
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMixedWrongSigMsg(ks.Shares[1], 1, version, ks))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMixedMsg(ks.Shares[2], 2, version, ks))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMixedMsg(ks.Shares[3], 3, version, ks))),
	//			},
	//			DontStartDuty:     true,
	//			ExpectedErrorCode: expectedErrorCode,
	//		},
	//	}...)
	//}

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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongSigAttestationMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongSigSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongSigAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[3], 3, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
			},
		}...)
	}

	return multiSpecTest
}
