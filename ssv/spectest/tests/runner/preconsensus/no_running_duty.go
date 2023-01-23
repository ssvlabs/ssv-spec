package preconsensus

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
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
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "ea0077ca087bf72425d9eefad026631277cf186dc5fc1dc71c3f5beab1395a4a",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing sync committee selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "19a9e1d4ff370a96c5605758264fd97ccf3d10ca863c5015ad6f8fe766b2a821",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "8a2416a0f62d227c0e823de5b5c54d859230376a2f7657eb1a71c859b3e2665c",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "a0b7a94f24ef679af43f3b41eff747ede38b9de0f693f8cc0025e756e0af8a4d",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "efcf694556d4ddeb6fb54b794b8cf04277e1a47ab161e827d48d6b1fc8d383e9",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing validator registration message: invalid pre-consensus message: no running duty",
			},
		},
	}
}
