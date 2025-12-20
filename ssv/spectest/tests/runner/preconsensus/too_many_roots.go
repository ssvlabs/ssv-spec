package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TooManyRoots tests too many expected roots
func TooManyRoots() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus too many roots",
		testdoc.PreConsensusTooManyRootsDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofTooManyRootsMsg(ks.Shares[1], ks.Shares[1], 1, 1, testingutils.TestingDutySlot))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				// No error is expected as AggregatorCommitteeRunner does not validate the number of roots
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoTooManyRootsMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.WrongRootsCountErrorCode,
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoTooManyRootsMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.WrongRootsCountErrorCode,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationTooManyRootsMsg(ks.Shares[1], 1))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.WrongRootsCountErrorCode,
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitTooManyRootsMsg(ks.Shares[1], 1))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.WrongRootsCountErrorCode,
			},
		},
		ks,
	)

	// Aggregator Committee duty
	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name:   "sync committee aggregator selection proof",
		Runner: testingutils.AggregatorCommitteeRunner(ks),
		Duty:   testingutils.TestingSyncCommitteeContributionDuty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofTooManyRootsMsg(ks.Shares[1], ks.Shares[1], 1, 1, testingutils.TestingDutySlot))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
		},
		// No error is expected as AggregatorCommitteeRunner does not validate the number of roots
	})
	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:   fmt.Sprintf("aggregator selection proof (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofTooManyRootsMsg(ks.Shares[1], ks.Shares[1], 1, 1, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version), // broadcasts when starting a new duty
				},
				// No error is expected as AggregatorCommitteeRunner does not validate the number of roots
			},
			{
				Name:   fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingAggregatorCommitteeDutyMixed(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsgTooManyRoots(ks.Shares[1], 1, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[1], 1, version), // broadcasts when starting a new duty
				},
				// No error is expected as AggregatorCommitteeRunner does not validate the number of roots
			},
		}...)
	}

	return multiSpecTest
}
