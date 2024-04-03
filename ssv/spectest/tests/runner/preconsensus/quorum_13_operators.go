package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Quorum13Operators  tests a quorum of valid PartialSignatureMessages (13 operators)
func Quorum13Operators() tests.SpecTest {
	ks := testingutils.Testing13SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus quorum 13 operators",
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
					testingutils.SSVMsgSyncCommitteeContribution(9, ks.NetworkKeys[9], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[9], 9)),
				},
				PostDutyRunnerStateRoot: "35c9e097c8559405a0ab316b162a63be3e004934c556c7c9d3699ceb77e43c92",
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
					testingutils.SSVMsgAggregator(9, ks.NetworkKeys[9], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[9], 9)),
				},
				PostDutyRunnerStateRoot: "23c970adc0e86c809bd15639eea2b6cd4e467aa51d473b13a71822594cf2910b",
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
					testingutils.SSVMsgProposer(9, ks.NetworkKeys[9], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[9], 9, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "b47777269ce127fa0446db86bcb8c261cf370b9c001bdf9d3b2b58a64de6aced",
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
					testingutils.SSVMsgProposer(9, ks.NetworkKeys[9], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[9], 9, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "9158e10db9cd8d86d17ac4edde8af4c837d3f193036bf4140badb8bc485f5842",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgValidatorRegistration(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgValidatorRegistration(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgValidatorRegistration(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgValidatorRegistration(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgValidatorRegistration(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[5], 5)),
					testingutils.SSVMsgValidatorRegistration(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[6], 6)),
					testingutils.SSVMsgValidatorRegistration(7, ks.NetworkKeys[7], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[7], 7)),
					testingutils.SSVMsgValidatorRegistration(8, ks.NetworkKeys[8], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[8], 8)),
					testingutils.SSVMsgValidatorRegistration(9, ks.NetworkKeys[9], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[9], 9)),
				},
				PostDutyRunnerStateRoot: "530111e289234d0c522eee11f16287de385592c2677f0acb18ba6435544738aa",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingValidatorRegistration),
				},
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgVoluntaryExit(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgVoluntaryExit(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgVoluntaryExit(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgVoluntaryExit(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgVoluntaryExit(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[5], 5)),
					testingutils.SSVMsgVoluntaryExit(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[6], 6)),
					testingutils.SSVMsgVoluntaryExit(7, ks.NetworkKeys[7], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[7], 7)),
					testingutils.SSVMsgVoluntaryExit(8, ks.NetworkKeys[8], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[8], 8)),
					testingutils.SSVMsgVoluntaryExit(9, ks.NetworkKeys[9], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[9], 9)),
				},
				PostDutyRunnerStateRoot: "530111e289234d0c522eee11f16287de385592c2677f0acb18ba6435544738aa",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedVoluntaryExit(ks)),
				},
			},
		},
	}
}
