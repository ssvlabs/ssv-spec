package postconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostFinish  tests a valid SignedPartialSignatureMessage post finished runner
func PostFinish() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	err := "failed processing post consensus message: invalid post-consensus message: no running duty"
	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: finishRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[4], 4, ks)),
				},
				PostDutyRunnerStateRoot: "197ca8360c7833627c664fc762512716a34b762acb85fd3e2467aae7b992e253",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name: "sync committee",
				Runner: finishRunner(
					testingutils.SyncCommitteeRunner(ks),
					&testingutils.TestingSyncCommitteeDuty,
					testingutils.TestSyncCommitteeConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[4], 4)),
				},
				PostDutyRunnerStateRoot: "378d58eb8d2e934c1eee08ee7c5a3cbda50db524b7b0d19ac8fcfdd2883e2b91",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name: "proposer",
				Runner: finishRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
					testingutils.TestProposerConsensusDataV(spec.DataVersionBellatrix),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, spec.DataVersionBellatrix)),
				},
				PostDutyRunnerStateRoot: "88a1ca7d2f30fe6714a452aec4f9bd4801ab504e32107784e1f30e872e17991d",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name: "proposer (blinded block)",
				Runner: finishRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
					testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionBellatrix),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, spec.DataVersionBellatrix)),
				},
				PostDutyRunnerStateRoot: "a647fb73237fc2bf0b4c690b03b311a7cfb0d7a2f8c22fcff02a41a4d55d93d7",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name: "aggregator",
				Runner: finishRunner(
					testingutils.AggregatorRunner(ks),
					&testingutils.TestingAggregatorDuty,
					testingutils.TestAggregatorConsensusData,
				),
				Duty: &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[4], 4)),
				},
				PostDutyRunnerStateRoot: "764d8da1b5c066bffa65ff8d534a430babccd897e7a1f8adb18d8988bd57208f",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name: "attester",
				Runner: finishRunner(
					testingutils.AttesterRunner(ks),
					&testingutils.TestingAttesterDuty,
					testingutils.TestAttesterConsensusData,
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[4], 4, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "f64ae938b2a3930ed0cfee30e40a75468403717afcbdf16a581fe884cf165306",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
		},
	}
}
