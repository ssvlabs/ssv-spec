package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidMessage tests a valid SignedPartialSignatureMessage with multi PartialSignatureMessages
func ValidMessage() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus valid msg",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "85036a9bf2d0b8ecaec3c539cc977efbfcaadde3a232c35c5ab4c1db45c7923d",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "7e90fccb1d93995a0079f522c4c8faab85fd0b260ab6051b90fa8befb8b0c98a",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "717e53a8cdb20f1d676bef0238637e42687c594758586f28d166697660ceb72d",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PreConsensusFailedMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "cd9cbb0cc68605e8b90f60c35d40558d23f6c709f53476357be83e1a46234d95",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "no pre consensus sigs required for attester role",
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PreConsensusFailedMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "8eacf6f368d395ae1c0d6b016963c74ad843958c31efe42c8661ffd3c337f106",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "no pre consensus sigs required for sync committee role",
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "ff6b65005d3abe4844d7a11b947d4e920da075ee47bac92e45584ca7ed875fdb",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
	}
}
