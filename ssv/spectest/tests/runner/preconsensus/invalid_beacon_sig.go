package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidBeaconSignature tests PartialSignatureMessage signature invalid
func InvalidBeaconSignature() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus invalid beacon signature",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofWrongBeaconSigMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "a5c4aab5b21e68e24be242ecebb90d4e711a0a77397e40dfbb6e7c67b5ff9d11",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing sync committee selection proof message: invalid pre-consensus message: could not verify Beacon partial Signature: wrong signature",
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofWrongBeaconSigMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "2ced65306471d08e0de69c2016d90c0651df26b3d5e231f12944704a48f59339",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing selection proof message: invalid pre-consensus message: could not verify Beacon partial Signature: wrong signature",
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoWrongBeaconSigMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "d46ee9acb3d6b3386c80e3947e83ac73063be1607b7fdd4a011f9814978082c4",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing randao message: invalid pre-consensus message: could not verify Beacon partial Signature: wrong signature",
			},
		},
	}
}
