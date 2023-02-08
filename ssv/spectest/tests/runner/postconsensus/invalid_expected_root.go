package postconsensus

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// InvalidExpectedRoot tests 1 expected root which doesn't match the signed root
func InvalidExpectedRoot() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedError := "failed processing post consensus message: invalid post-consensus message: wrong signing root"
	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus invalid expected roots",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: decideRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
				),
				Duty: testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusWrongSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
				},
				PostDutyRunnerStateRoot: "5f70323f99cc96693e30d0ebd9ec3dac7f75948045e4089c7ba630f6c0b5422d",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "sync committee",
				Runner: decideRunner(
					testingutils.SyncCommitteeRunner(ks),
					testingutils.TestingSyncCommitteeDuty,
					testingutils.TestSyncCommitteeConsensusData,
				),
				Duty: testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusWrongSyncCommitteeMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "eff09708ecaa0cd6e19ad8e841592a13eda84da98526be66eaac292e9d340fdc",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "proposer",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDuty,
					testingutils.TestProposerConsensusData,
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongProposerMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "5ffdf605dded2c2e140cbca6ef8b0cb1b4192cf087530bc9d36f5a2f907c6b84",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "proposer (blinded block)",
				Runner: decideRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestingProposerDuty,
					testingutils.TestProposerBlindedBlockConsensusData,
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongProposerMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "a2e539418ce3768e21786ab7a7377a100d6064d6a02fb6c0b83aad2f7e12485c",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "aggregator",
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					testingutils.TestingAggregatorDuty,
					testingutils.TestAggregatorConsensusData,
				),
				Duty: testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusWrongAggregatorMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "177d2430aafdb3efadf58e342d60e5e324adf7798cea07744d45b6ea68e97245",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "attester",
				Runner: decideRunner(
					testingutils.AttesterRunner(ks),
					testingutils.TestingAttesterDuty,
					testingutils.TestAttesterConsensusData,
				),
				Duty: testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusWrongAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "2741163569f537af394880fa208fbb916cfea41280403a1ecd509e26d9f98d4d",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
		},
	}
}
