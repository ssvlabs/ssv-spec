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
				PostDutyRunnerStateRoot: "505247fc7a3303acc63b7863d1d05567de2878cc470aeb221b898758ad040ffc",
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
				PostDutyRunnerStateRoot: "3f416a6c4d34ceadecdfbad133f5ef4880ffb2915f962804d073f2b2a1a49fe2",
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
				PostDutyRunnerStateRoot: "2c86b28ed450ee7a356ab5cdd330c303297da5b0ca5e50e044c9350e88fc004b",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name: "randao (blinded block)",
				Runner: decideRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestingProposerDuty,
					testingutils.TestProposerBlindedBlockConsensusData,
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
				PostDutyRunnerStateRoot: "50593bc477c0a77cc07186cda738a06c57389dc4a57aba3e0283a768545a4f27",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
