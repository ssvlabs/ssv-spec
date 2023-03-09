package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoRunningDuty tests a valid partial pre consensus msg before duty starts
func NoRunningDuty() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus no running duty",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "843d535c51be1633f4d48589e339b8f7ff38185e094292fb7a7117f53c452019",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing sync committee selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "f5aec3ffac19854712478ecd8a97f2b1676c35bcf840d13da0d18aebfa8097d5",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "4e23bc72477fcde5fd2301bb20f5b12303fe1516c5c238d319cf108cc98b8153",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "7e3497d83dfb9d984a1138c2118150e9f6041b03e936d7c42c6cfcf980c03d38",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "194ee4fd2a60a8c1f4d4d51f4445e4176e199ac075d06614a121c0d13fc63e0b",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing validator registration message: invalid pre-consensus message: no running duty",
			},
		},
	}
}
