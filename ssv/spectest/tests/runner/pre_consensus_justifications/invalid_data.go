package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidData tests error decoding consensusData
func InvalidData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// invalidateMsgDataF sets a non ValidatorConsensusData data in message, it will fail when decoding
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

	expectedErr := "failed processing consensus message: invalid pre-consensus justification: could not decoded ValidatorConsensusData: incorrect offset"

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
		},
	}
}
