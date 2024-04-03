package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidMessage13Operators tests a valid PartialSignatureMessages with multi PartialSignatureMessages (13 operators)
func ValidMessage13Operators() tests.SpecTest {
	ks := testingutils.Testing13SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus valid msg 13 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgSyncCommitteeContribution(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommitteeContribution(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgSyncCommitteeContribution(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgSyncCommitteeContribution(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[5], 5)),
					testingutils.SSVMsgSyncCommitteeContribution(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[6], 6)),
					testingutils.SSVMsgSyncCommitteeContribution(7, ks.NetworkKeys[7], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[7], 7)),
					testingutils.SSVMsgSyncCommitteeContribution(8, ks.NetworkKeys[8], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[8], 8)),
				},
				PostDutyRunnerStateRoot: "f3f3a82eb2a42fc4d76fc1a6ce922ddf55948f7d01c091127401aec25c71d401",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgAggregator(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgAggregator(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgAggregator(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgAggregator(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgAggregator(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[5], 5)),
					testingutils.SSVMsgAggregator(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[6], 6)),
					testingutils.SSVMsgAggregator(7, ks.NetworkKeys[7], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[7], 7)),
					testingutils.SSVMsgAggregator(8, ks.NetworkKeys[8], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[8], 8)),
				},
				PostDutyRunnerStateRoot: "43862817ab29a13e9445d32a6d39ae31d8652b8855fbff3e3db1f24805ecbf8a",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], 5, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[6], 6, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(7, ks.NetworkKeys[7], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[7], 7, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(8, ks.NetworkKeys[8], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[8], 8, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "11f25e657369c00dcde953ef3a0f732e774c716f8089f7e624fd4409331dce05",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], 5, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[6], 6, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(7, ks.NetworkKeys[7], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[7], 7, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(8, ks.NetworkKeys[8], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[8], 8, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "62da0545fcfea1f330db88e07dcb6f9391903bc6ee6a34d23ae6345f355ae505",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
		},
	}
}
