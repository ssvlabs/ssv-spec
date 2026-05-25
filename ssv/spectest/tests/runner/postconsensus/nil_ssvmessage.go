package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NilSSVMessage tests a SignedSSVMessage with wrong data that can't be decoded
func NilSSVMessage() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	expectedErrorCode := types.NilSSVMessageErrorCode
	invalidMsg := &types.SignedSSVMessage{
		Signatures:  [][]byte{{1, 2, 3, 4}},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  nil,
	}

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus nil ssvmessage",
		testdoc.PostConsensusNilMsgDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name: "proposer",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb),
				),
				Duty:              testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages:          []*types.SignedSSVMessage{invalidMsg},
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
				Duty:              testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages:          []*types.SignedSSVMessage{invalidMsg},
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
		Duty:              testingutils.TestingSyncCommitteeContributionDuty,
		Messages:          []*types.SignedSSVMessage{invalidMsg},
		DontStartDuty:     true,
		ExpectedErrorCode: expectedErrorCode,
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
				Duty:              testingutils.TestingAggregatorDuty(version),
				Messages:          []*types.SignedSSVMessage{invalidMsg},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name: fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Runner: decideAggregatorCommitteeRunner(
					testingutils.AggregatorCommitteeRunner(ks),
					testingutils.TestingAggregatorCommitteeDutyMixed(version),
					testingutils.TestAggregatorCommitteeConsensusData(version),
				),
				Duty:              testingutils.TestingAggregatorCommitteeDutyMixed(version),
				Messages:          []*types.SignedSSVMessage{invalidMsg},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
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
				Duty:              testingutils.TestingAttesterDuty(version),
				Messages:          []*types.SignedSSVMessage{invalidMsg},
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
				Duty:              testingutils.TestingSyncCommitteeDuty(version),
				Messages:          []*types.SignedSSVMessage{invalidMsg},
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
				Duty:              testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages:          []*types.SignedSSVMessage{invalidMsg},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
			},
		}...)
	}

	return multiSpecTest
}
