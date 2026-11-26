package preconsensus

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NilSSVMessage tests a SignedSSVMessage with a nil SSVMessage
func NilSSVMessage() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedErrorCode := types.NilSSVMessageErrorCode

	invalidMsg := &types.SignedSSVMessage{
		Signatures:  [][]byte{{1, 2, 3, 4}},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  nil,
	}

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus nil ssvmessage",
		testdoc.PreConsensusNilMsgDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:                    "randao",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(gloas.DataVersionGloas),
				Messages:                []*types.SignedSSVMessage{invalidMsg},
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
		Name:                    "sync committee aggregator selection proof",
		Runner:                  testingutils.AggregatorCommitteeRunner(ks),
		Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
		Messages:                []*types.SignedSSVMessage{invalidMsg},
		PostDutyRunnerStateRoot: "29862cc6054edc8547efcb5ae753290971d664b9c39768503b4d66e1b52ecb06",
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
		},
		ExpectedErrorCode: expectedErrorCode,
	})
	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:     fmt.Sprintf("aggregator selection proof (%s)", version.String()),
				Runner:   testingutils.AggregatorCommitteeRunner(ks),
				Duty:     testingutils.TestingAggregatorDuty(version),
				Messages: []*types.SignedSSVMessage{invalidMsg},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:     fmt.Sprintf("aggregator committee duty (%s)", version.String()),
				Runner:   testingutils.AggregatorCommitteeRunner(ks),
				Duty:     testingutils.TestingAggregatorCommitteeDutyMixed(version),
				Messages: []*types.SignedSSVMessage{invalidMsg},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[1], 1, version), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
		}...)
	}

	return multiSpecTest
}
