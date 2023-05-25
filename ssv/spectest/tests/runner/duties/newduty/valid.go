package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Valid tests a valid start duty
func Valid() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty valid",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "468dc1442b1c5b845ef246e95cc0b62c0df32eb70f18a716311a7d247cb6b668",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "5f169237efd54ad718c306a184ddf79e26e781cb18345c8128653fd591b2d0f0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    &testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "6a03b13e954bc24c7f3efe6f15e2b79a80befef5120dd556815f96f7e04d7d6e",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "ff672c80ed3895981e225b94921d3cce7f1518d807a96afb964bbba9bb596d9a",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "1116d498c93294fc5b828450f1dd2f20409395fc0dd0ace9b22bf1ef61dff82a",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
