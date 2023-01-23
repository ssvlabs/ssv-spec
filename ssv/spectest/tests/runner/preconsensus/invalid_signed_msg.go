package preconsensus

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// InvalidSignedMessage tests SignedPartialSignatureMessage.Validate() != nil
func InvalidSignedMessage() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	invaldiateMsg := func(msg *ssv.SignedPartialSignatureMessage) *ssv.SignedPartialSignatureMessage {
		msg.Signature = nil
		return msg
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus invalid signed msg",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, invaldiateMsg(testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
				},
				PostDutyRunnerStateRoot: "93f63d09f9a12244c8ca751934450e7be81f89eb876dc76119ae7a7153ed405c",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing sync committee selection proof message: invalid pre-consensus message: SignedPartialSignatureMessage invalid: SignedPartialSignatureMessage sig invalid",
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, invaldiateMsg(testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
				},
				PostDutyRunnerStateRoot: "57a6fb4d6bf799353c3c0980e3eb31ff32f0c100ebe04f6025dc06fc983ef52b",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing selection proof message: invalid pre-consensus message: SignedPartialSignatureMessage invalid: SignedPartialSignatureMessage sig invalid",
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, invaldiateMsg(testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
				},
				PostDutyRunnerStateRoot: "65cd178a17356e48a1d7de571c0d56f791f6c4654b5b03e385adea510f7c55d1",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing randao message: invalid pre-consensus message: SignedPartialSignatureMessage invalid: SignedPartialSignatureMessage sig invalid",
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, invaldiateMsg(testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
				},
				PostDutyRunnerStateRoot: "2a122d8afb55f8cf02b7008d9e525d4dbd5dd839e752be9dd5be577e653c56e4",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing randao message: invalid pre-consensus message: SignedPartialSignatureMessage invalid: SignedPartialSignatureMessage sig invalid",
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(nil, invaldiateMsg(testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: "34b8788e99ec86321d758e7a55a7e688bb91cc15d0a8fd581cf28f0839601fbd",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing validator registration message: invalid pre-consensus message: SignedPartialSignatureMessage invalid: SignedPartialSignatureMessage sig invalid",
			},
		},
	}
}
