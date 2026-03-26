package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TooFewRoots tests a valid SignedPartialSignatureMessage with too few roots
func TooFewRoots() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	errCode := types.NoPartialSigMessagesErrorCode
	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus too few roots",
		testdoc.PostConsensusTooFewRootsDoc,
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
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerTooFewRootsMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
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
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerTooFewRootsMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
			},
		},
		ks,
	)

	// Aggregator committee duty
	sccSlot := testingutils.TestingSyncCommitteeContributionDuty.Slot
	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name: "sync committee contribution",
		Runner: decideAggregatorCommitteeRunner(
			testingutils.AggregatorCommitteeRunner(ks),
			testingutils.TestingSyncCommitteeContributionDuty,
			testingutils.TestSyncCommitteeContributionConsensusData,
		),
		Duty: testingutils.TestingSyncCommitteeContributionDuty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionTooFewRootsMsg(ks.Shares[1], 1, ks, sccSlot))),
		},
		DontStartDuty: true,
		// No error is expected as AggregatorCommitteeRunner doesn't validate the precise number of roots
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
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorTooFewRootsMsg(ks.Shares[1], 1, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
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
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMixedMsgTooFewRoots(ks.Shares[1], 1, version, ks))),
				},
				DontStartDuty: true,
				// No error is expected as AggregatorCommitteeRunner doesn't validate the precise number of roots
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationTooFewRootsMsg(ks.Shares[1], 1, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeTooFewRootsMsg(ks.Shares[1], 1, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgTooFewRootsMsg(ks.Shares[1], 1, version))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: errCode,
			},
		}...)
	}

	return multiSpecTest
}
