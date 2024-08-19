package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownSigner tests SignedPartialSignatureMessage unknown signer
func UnknownSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	expectedError := "SignedSSVMessage has an invalid signature: unknown signer"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus unknown msg signer",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageWithSigner(5, ks.OperatorKeys[1], testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 5, 5))),
				},
				PostDutyRunnerStateRoot: unknownSignerSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     unknownSignerSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageWithSigner(5, ks.OperatorKeys[1], testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 5, 5))),
				},
				PostDutyRunnerStateRoot: unknownSignerAggregatorSC().Root(),
				PostDutyRunnerState:     unknownSignerAggregatorSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("randao (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignedSSVMessageWithSigner(5, ks.OperatorKeys[1], testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 5, 5, version))),
			},
			PostDutyRunnerStateRoot: unknownSignerProposerSC(version).Root(),
			PostDutyRunnerState:     unknownSignerProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
			ExpectedError: expectedError,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("randao blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignedSSVMessageWithSigner(5, ks.OperatorKeys[1], testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 5, 5, version))),
			},
			PostDutyRunnerStateRoot: unknownSignerBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     unknownSignerBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
			ExpectedError: expectedError,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
