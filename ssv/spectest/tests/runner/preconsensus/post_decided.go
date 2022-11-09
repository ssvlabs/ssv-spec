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
	decideRunner := func(r ssv.Runner, duty *types.Duty, decidedValue *types.ConsensusData, preMsgs []*ssv.SignedPartialSignatures) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		for _, msg := range preMsgs {
			r.ProcessPreConsensus(msg)
		}
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
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
					[]*ssv.SignedPartialSignatures{
						testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
						testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2),
						testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3),
					},
				),
				Duty: testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[4], ks.Shares[4], 4, 4), types.PartialContributionProofSignatureMsgType),
				},
				PostDutyRunnerStateRoot: "dc06fe49582c438ece971da8aa8107e4ed90ca1000403d43aa64564e4be1e9d2",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatures{},
			},
			{
				Name: "aggregator selection proof",
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					testingutils.TestingAggregatorDuty,
					testingutils.TestAggregatorConsensusData,
					[]*ssv.SignedPartialSignatures{
						testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
						testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2),
						testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3),
					},
				),
				Duty: testingutils.TestingAggregatorDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], ks.Shares[4], 4, 4), types.PartialSelectionProofSignatureMsgType),
				},
				PostDutyRunnerStateRoot: "bb1378ee14f4e8bb55880fe42d41eb849276aa71533b05c777f4783b115f0d3a",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatures{},
			},
			{
				Name: "randao",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDuty,
					testingutils.TestProposerConsensusData,
					[]*ssv.SignedPartialSignatures{
						testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1),
						testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[2], ks.Shares[2], 2, 2),
						testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[3], ks.Shares[3], 3, 3),
					},
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[4], ks.Shares[4], 4, 4), types.PartialRandaoSignatureMsgType),
				},
				PostDutyRunnerStateRoot: "4b7c69ae85975112b30a3cb233319e3102bdef7cd3962d498d90d74b6a50a016",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatures{},
			},
		},
	}
}
