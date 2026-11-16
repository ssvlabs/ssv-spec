package preconsensus

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidBeaconSignatureInQuorum tests a quorum with an invalid PartialSignatureMessage signature
func InvalidBeaconSignatureInQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedErrorCode := types.ReconstructSignatureErrorCode
	sccSlot := testingutils.TestingSyncCommitteeContributionDuty.Slot
	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus invalid beacon signature in quorum",
		testdoc.PreConsensusInvalidBeaconSignatureInQuorumDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(gloas.DataVersionGloas),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoWrongBeaconSigMsgV(ks.Shares[1], 1, gloas.DataVersionGloas))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, gloas.DataVersionGloas))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, gloas.DataVersionGloas))),
				},
				PostDutyRunnerStateRoot: "56eafcb33392ded888a0fefe30ba49e52aa00ab36841cb10c9dc1aa2935af347",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, gloas.DataVersionGloas), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
		},
		ks,
	)

	// Aggregator committee duty
	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name:   "sync committee aggregator selection proof",
		Runner: testingutils.AggregatorCommitteeRunner(ks),
		Duty:   testingutils.TestingSyncCommitteeContributionDuty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofWrongBeaconSigMsg(ks.Shares[1], ks.Shares[1], 1, 1, sccSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[2], ks.Shares[2], 2, 2, sccSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[3], ks.Shares[3], 3, 3, sccSlot))),
		},
		PostDutyRunnerStateRoot: "29862cc6054edc8547efcb5ae753290971d664b9c39768503b4d66e1b52ecb06",
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[1], ks.Shares[1], 1, 1, sccSlot), // broadcasts when starting a new duty
		},
		ExpectedErrorCode: expectedErrorCode,
	})
	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:   fmt.Sprintf("aggregator selection proof (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofWrongBeaconSigMsg(ks.Shares[1], ks.Shares[1], 1, 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingAggregatorCommitteeDutyMixed(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsgWrongBeaconSig(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[3], 3, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[1], 1, version), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
		}...)
	}

	return multiSpecTest
}
