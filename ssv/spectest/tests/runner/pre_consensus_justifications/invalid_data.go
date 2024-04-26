package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidData tests error decoding consensusData
func InvalidData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// invalidateMsgDataF sets a non ConsensusData data in message, it will fail when decoding
	invalidateMsgDataF := func(id []byte) *types.SignedSSVMessage {
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     1,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       testingutils.TestingQBFTRootData,
		}
		signed := testingutils.SignQBFTMsg(ks.OperatorKeys[1], 1, msg)
		signed.FullData = testingutils.TestingQBFTFullData

		return signed
	}

	expectedErr := "failed processing consensus message: invalid pre-consensus justification: could not decoded ConsensusData: incorrect offset"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus invalid data",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: decideFirstHeight(testingutils.SyncCommitteeContributionRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					invalidateMsgDataF(testingutils.SyncCommitteeContributionMsgID),
				},
				PostDutyRunnerStateRoot: "2619aeecde47fe0efc36aa98fbb2df9834d9eee77f62abe0d10532dbd5215790",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "aggregator selection proof",
				Runner: decideFirstHeight(testingutils.AggregatorRunner(ks)),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					invalidateMsgDataF(testingutils.AggregatorMsgID),
				},
				PostDutyRunnerStateRoot: "db1b416873d19be76cddc92ded0d442ba0e642514973b5dfec45f587c6ffde15",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "randao",
				Runner: decideFirstHeight(testingutils.ProposerRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					invalidateMsgDataF(testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "2754fc7ced14fb15f3f18556bb6b837620287cbbfbf908abafa5a0533fc4bc5f",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "randao (blinded block)",
				Runner: decideFirstHeight(testingutils.ProposerBlindedBlockRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					invalidateMsgDataF(testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "6bd59da9f817b8e40112e58231e36738b9d021db4416c9eeec1dd0236a5362e2",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			// {

			// 	Name:   "attester",
			// 	Runner: decideFirstHeight(testingutils.CommitteeRunner(ks)),
			// 	Duty:   &testingutils.TestingAttesterDuty,
			// 	Messages: []*types.SignedSSVMessage{
			// 		invalidateMsgDataF(testingutils.AttesterMsgID),
			// 	},
			// 	PostDutyRunnerStateRoot: "81cb7b1d3ea3087d49f9773b3a2b75a87b901e50427d237f2a10c0e1904e7684",
			// 	OutputMessages:          []*types.PartialSignatureMessages{},
			// 	ExpectedError:           "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: failed decoding consensus data: incorrect offset",
			// },
			// {
			// 	Name:   "sync committee",
			// 	Runner: decideFirstHeight(testingutils.SyncCommitteeRunner(ks)),
			// 	Duty:   &testingutils.TestingSyncCommitteeDuty,
			// 	Messages: []*types.SignedSSVMessage{
			// 		invalidateMsgDataF(testingutils.SyncCommitteeMsgID),
			// 	},
			// 	PostDutyRunnerStateRoot: "38592232077cd45709a7c6cfdd20c9d899af9d79bc750add3c4b8f2b6794cb34",
			// 	OutputMessages:          []*types.PartialSignatureMessages{},
			// 	ExpectedError:           "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: failed decoding consensus data: incorrect offset",
			// },
		},
	}
}
