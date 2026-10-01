package preconsensus

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidExpectedRoot tests 1 expected root which doesn't match the signed root
func InvalidExpectedRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus invalid expected roots",
		testdoc.PreConsensusInvalidExpectedRootDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(gloas.DataVersionGloas),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentEpochMsgV(ks.Shares[1], 1, gloas.DataVersionGloas))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, gloas.DataVersionGloas), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.WrongSigningRootErrorCode,
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitWrongRootMsg(ks.Shares[1], 1))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.WrongSigningRootErrorCode,
			},
			{
				Name:   "ptc attestation",
				Runner: testingutils.PTCAttesterRunner(ks),
				Duty:   testingutils.TestingPTCAttesterDuty(),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPTCAttester(nil, testingutils.PreConsensusPTCWrongRootMsg(ks.Shares[1], 1))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusPTCMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.WrongSigningRootErrorCode,
			},
			{
				Name:   "proposer preferences",
				Runner: testingutils.ProposerPreferencesRunner(ks),
				Duty:   testingutils.TestingProposerPreferencesDuty(),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesWrongRootMsg(ks.Shares[1], 1))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.WrongSigningRootErrorCode,
			},
		},
		ks,
	)

	// Aggregator Committee duty
	sccSlot := testingutils.TestingSyncCommitteeContributionDuty.Slot
	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name:   "sync committee aggregator selection proof",
		Runner: testingutils.AggregatorCommitteeRunner(ks),
		Duty:   testingutils.TestingSyncCommitteeContributionDuty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofWrongBeaconRootMsg(ks.Shares[1], ks.Shares[1], 1, 1, sccSlot))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[1], ks.Shares[1], 1, 1, sccSlot), // broadcasts when starting a new duty
		},
		// No error as AggregatorCommitteeRunner does not validate the expected root
	})
	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:   fmt.Sprintf("aggregator selection proof (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofWrongRootSigMsg(ks.Shares[1], ks.Shares[1], 1, 1, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version), // broadcasts when starting a new duty
				},
				// No error as AggregatorCommitteeRunner does not validate the expected root
			},
			{
				Name:   fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingAggregatorCommitteeDutyMixed(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsgWrongRoot(ks.Shares[1], 1, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[1], 1, version), // broadcasts when starting a new duty
				},
				// No error as AggregatorCommitteeRunner does not validate the expected root
			},
		}...)
	}

	return multiSpecTest
}
