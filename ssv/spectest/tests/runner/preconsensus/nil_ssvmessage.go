package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NilSSVMessage tests a SignedSSVMessage with a nil SSVMessage
func NilSSVMessage() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedError := "invalid SignedSSVMessage: nil SSVMessage"

	invalidMsg := &types.SignedSSVMessage{
		Signatures:  [][]byte{{1, 2, 3, 4}},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  nil,
	}

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus nil ssvmessage",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee aggregator selection proof",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                []*types.SignedSSVMessage{invalidMsg},
				PostDutyRunnerStateRoot: "29862cc6054edc8547efcb5ae753290971d664b9c39768503b4d66e1b52ecb06",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "randao",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages:                []*types.SignedSSVMessage{invalidMsg},
				PostDutyRunnerStateRoot: "56eafcb33392ded888a0fefe30ba49e52aa00ab36841cb10c9dc1aa2935af347",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "randao (blinded block)",
				Runner:                  testingutils.ProposerBlindedBlockRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages:                []*types.SignedSSVMessage{invalidMsg},
				PostDutyRunnerStateRoot: "2ce3241658f324f352c77909f4043934eedf38e939ae638c5ce6acf28e965646",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
		},
	}

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:     fmt.Sprintf("aggregator selection proof (%s)", version.String()),
			Runner:   testingutils.AggregatorRunner(ks),
			Duty:     testingutils.TestingAggregatorDuty(version),
			Messages: []*types.SignedSSVMessage{invalidMsg},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version), // broadcasts when starting a new duty
			},
			ExpectedError: expectedError,
		})
	}

	return multiSpecTest
}
