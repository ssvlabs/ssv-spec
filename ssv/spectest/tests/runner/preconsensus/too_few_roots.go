package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TooFewRoots tests too few expected roots
func TooFewRoots() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedErrorCode := types.NoPartialSigMessagesErrorCode
	sccSlot := testingutils.TestingSyncCommitteeContributionDuty.Slot
	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus too few roots",
		testdoc.PreConsensusTooFewRootsDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofTooFewRootsMsg(ks.Shares[1], ks.Shares[1], 1, 1, sccSlot))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				// AggregatorCommitteeRunner doesn't validate the number of roots for SC
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoTooFewRootsMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoTooFewRootsMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationTooFewRootsMsg(ks.Shares[1], 1))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitTooFewRootsMsg(ks.Shares[1], 1))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
		},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("aggregator selection proof (%s)", version.String()),
			Runner: testingutils.AggregatorCommitteeRunner(ks),
			Duty:   testingutils.TestingAggregatorDuty(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofTooFewRootsMsg(ks.Shares[1], ks.Shares[1], 1, 1, version))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version), // broadcasts when starting a new duty
			},
			ExpectedErrorCode: expectedErrorCode,
		})
	}

	return multiSpecTest
}
