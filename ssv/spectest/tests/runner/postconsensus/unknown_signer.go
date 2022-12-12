package postconsensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests an unknown  signer
func UnknownSigner() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedError := "failed processing post consensus message: invalid post-consensus message: failed to verify PartialSignature: unknown signer"
	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus unknown signer",
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
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSigSyncCommitteeContributionWrongBeaconSignerMsg(ks.Shares[1], 5, 1, ks)),
				},
				PostDutyRunnerStateRoot: "d12a562cba23fe156380cb61200ab8abe6aec9dd90a6842040b2ef57e50f26a2",
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
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSigSyncCommitteeWrongBeaconSignerMsg(ks.Shares[1], 5, 1)),
				},
				PostDutyRunnerStateRoot: "4e84eca5e18d0aebcac7de43d886c04c3182c37d00f13e1040a80b215bce918e",
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
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusSigProposerWrongBeaconSignerMsg(ks.Shares[1], 5, 1)),
				},
				PostDutyRunnerStateRoot: "8750d044bb1d5d919c901a889a7058241744c89b8a5df42714a4af091bfa387c",
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
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusSigAggregatorWrongBeaconSignerMsg(ks.Shares[1], 5, 1)),
				},
				PostDutyRunnerStateRoot: "59756858e7fb953f3273ba26ea4e48b491eebc98abb2f86054bd1d62d7a1d620",
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
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusSigAttestationWrongBeaconSignerMsg(ks.Shares[1], 5, 1, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "09d81b254934d9bd37b28d3e72ffc5e0fce2d3a0c61b3e4204b9ba4043270869",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
		},
	}
}
