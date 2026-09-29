package preconsensus

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidSignedMessageDifferentLength tests an invalid SignedSSVMessage with a different number of signers and signatures
func InvalidSignedMessageDifferentLength() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	diffLength := func(msg *types.SignedSSVMessage) *types.SignedSSVMessage {
		msg.Signatures = append(msg.Signatures, []byte{1, 2, 3, 4})
		return msg
	}

	expectedErrorCode := types.IncorrectNumberOfSignaturesErrorCode

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus invalid signed msg different length",
		testdoc.PreConsensusInvalidSignedMessageDifferentLengthDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(gloas.DataVersionGloas),
				Messages: []*types.SignedSSVMessage{
					diffLength(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, gloas.DataVersionGloas)))),
				},
				PostDutyRunnerStateRoot: "56eafcb33392ded888a0fefe30ba49e52aa00ab36841cb10c9dc1aa2935af347",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, gloas.DataVersionGloas), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					diffLength(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1)))),
				},
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   "ptc attestation",
				Runner: testingutils.PTCAttesterRunner(ks),
				Duty:   testingutils.TestingPTCAttesterDuty(),
				Messages: []*types.SignedSSVMessage{
					diffLength(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgPTCAttester(nil, testingutils.PreConsensusPTCMsg(ks.Shares[1], 1)))),
				},
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusPTCMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   "proposer preferences",
				Runner: testingutils.ProposerPreferencesRunner(ks),
				Duty:   testingutils.TestingProposerPreferencesDuty(),
				Messages: []*types.SignedSSVMessage{
					diffLength(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1)))),
				},
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
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
			diffLength(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)))),
		},
		PostDutyRunnerStateRoot: "29862cc6054edc8547efcb5ae753290971d664b9c39768503b4d66e1b52ecb06",
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
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
					diffLength(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version)))),
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
					diffLength(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[1], 1, version)))),
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
