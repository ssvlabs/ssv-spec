package preconsensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests a msg received post consensus decided (and post receiving a quorum for pre consensus)
func PostDecided() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check errors
	// nolint
	decideRunner := func(r ssv.Runner, duty *types.Duty, decidedValue *types.ConsensusData, preMsgs []*ssv.SignedPartialSignatureMessage) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		for _, msg := range preMsgs {
			r.ProcessPreConsensus(msg)
		}
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GenerateConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().State.DecidedValue = decidedValue
		r.GetBaseRunner().QBFTController.StoredInstances[0] = r.GetBaseRunner().State.RunningInstance
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
		return r
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus post decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee aggregator selection proof",
				Runner: decideRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
					[]*ssv.SignedPartialSignatureMessage{
						testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
						testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2),
						testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3),
					},
				),
				Duty: testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[4], ks.Shares[4], 4, 4)),
				},
				PostDutyRunnerStateRoot: "366320a75151b43e701750fefe361405fb53e39269db1830314f881cac48126d",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name: "aggregator selection proof",
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					testingutils.TestingAggregatorDuty,
					testingutils.TestAggregatorConsensusData,
					[]*ssv.SignedPartialSignatureMessage{
						testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
						testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2),
						testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3),
					},
				),
				Duty: testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], ks.Shares[4], 4, 4)),
				},
				PostDutyRunnerStateRoot: "123efb0793124df79dcae660aa87918726a0dcde94031223394975e37ee0e098",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name: "randao",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDuty,
					testingutils.TestProposerConsensusData,
					[]*ssv.SignedPartialSignatureMessage{
						testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1),
						testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[2], ks.Shares[2], 2, 2),
						testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[3], ks.Shares[3], 3, 3),
					},
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[4], ks.Shares[4], 4, 4)),
				},
				PostDutyRunnerStateRoot: "a347ece3f05f70d34db60b25fe499ed82bc1e2c9c8eeef0dc680c9c8cd168c23",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
