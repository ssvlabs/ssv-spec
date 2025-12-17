package consensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidDecided7Operators tests a valid decided value (7 operators)
func ValidDecided7Operators() tests.SpecTest {

	ks := testingutils.Testing7SharesSet()

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"consensus valid decided 7 operators",
		testdoc.ConsensusValidDecided7OperatorsDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:     "sync committee contribution",
				Runner:   testingutils.AggregatorCommitteeRunner(ks),
				Duty:     testingutils.TestingSyncCommitteeContributionDuty,
				Messages: testingutils.SSVDecidingMsgsForAggregatorCommitteeRunnerForKS(testingutils.TestingSyncCommitteeContributionDuty, ks, spec.DataVersionPhase0),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
		},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:     fmt.Sprintf("aggregator (%s)", version.String()),
			Runner:   testingutils.AggregatorCommitteeRunner(ks),
			Duty:     testingutils.TestingAggregatorDuty(version),
			Messages: testingutils.SSVDecidingMsgsForAggregatorCommitteeRunnerForKS(testingutils.TestingAggregatorDuty(version), ks, version),
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

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:     fmt.Sprintf("proposer (%s)", version.String()),
			Runner:   testingutils.ProposerRunner(ks),
			Duty:     testingutils.TestingProposerDutyV(version),
			Messages: testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer),
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:     fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner:   testingutils.ProposerBlindedBlockRunner(ks),
			Duty:     testingutils.TestingProposerDutyV(version),
			Messages: testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(version), ks, types.RoleProposer),
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
