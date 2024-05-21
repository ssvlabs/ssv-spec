package consensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostDecided tests a valid commit msg after returned decided already
func PostDecided() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "attester",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, testingutils.TestingDutySlot),
					testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.CommitteeMsgID(ks), testingutils.TestBeaconVoteByts),
				),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, testingutils.TestingDutySlot),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, testingutils.TestingDutySlot),
					testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.CommitteeMsgID(ks), testingutils.TestBeaconVoteByts),
				),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "attester and sync committee",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterAndSyncCommitteeDuties,
				Messages: append(
					testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, testingutils.TestingDutySlot),
					testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.CommitteeMsgID(ks), testingutils.TestBeaconVoteByts),
				),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, testingutils.TestingDutySlot),
				},
			},
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.RoleSyncCommitteeContribution),
					testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.SyncCommitteeContributionMsgID, testingutils.TestSyncCommitteeContributionConsensusDataByts),
				),
				PostDutyRunnerStateRoot: postDecidedSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     postDecidedSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData, ks, types.RoleAggregator),
					testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.AggregatorMsgID, testingutils.TestAggregatorConsensusDataByts),
				),
				PostDutyRunnerStateRoot: postDecidedAggregatorSC().Root(),
				PostDutyRunnerState:     postDecidedAggregatorSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer),
				testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4],
					types.OperatorID(4), qbft.Height(testingutils.TestingDutySlotV(version)), testingutils.ProposerMsgID,
					testingutils.TestProposerConsensusDataBytsV(version)),
			),
			PostDutyRunnerStateRoot: postDecidedProposerSC(version).Root(),
			PostDutyRunnerState:     postDecidedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(version), ks, types.RoleProposer),
				testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4],
					types.OperatorID(4), qbft.Height(testingutils.TestingDutySlotV(version)), testingutils.ProposerMsgID,
					testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)),
			),
			PostDutyRunnerStateRoot: postDecidedBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     postDecidedBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
