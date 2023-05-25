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
	invalidateMsgDataF := func(id []byte) *qbft.SignedMessage {
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     1,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       testingutils.TestingQBFTRootData,
		}
		signed := testingutils.SignQBFTMsg(ks.Shares[1], 1, msg)
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
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(invalidateMsgDataF(testingutils.SyncCommitteeContributionMsgID), nil),
				},
				PostDutyRunnerStateRoot: "2e0d0c26372ecd5d3ff786bd28581b624fb2f30c84821b23bb15a61faf377d3e",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "aggregator selection proof",
				Runner: decideFirstHeight(testingutils.AggregatorRunner(ks)),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(invalidateMsgDataF(testingutils.AggregatorMsgID), nil),
				},
				PostDutyRunnerStateRoot: "c2433300980af6fae62151b9fc6bda67842e019205b5a31660692ab8e99ddbb2",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "randao",
				Runner: decideFirstHeight(testingutils.ProposerRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(invalidateMsgDataF(testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "3544989d49ef3071258fd72b2befa80366afafcb788914a17f6226ad957eb980",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "randao (blinded block)",
				Runner: decideFirstHeight(testingutils.ProposerBlindedBlockRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(invalidateMsgDataF(testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "35fce5d0dc47c393efcdaf98168a6b00e931695c73f6cc47f90d7d6b706ca089",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{

				Name:   "attester",
				Runner: decideFirstHeight(testingutils.AttesterRunner(ks)),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(invalidateMsgDataF(testingutils.AttesterMsgID), nil),
				},
				PostDutyRunnerStateRoot: "127cb9d3d0f38967765066de4a2512e2dc1f12f501b395afb0b6b697316e021c",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: failed decoding consensus data: incorrect offset",
			},
			{
				Name:   "sync committee",
				Runner: decideFirstHeight(testingutils.SyncCommitteeRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(invalidateMsgDataF(testingutils.SyncCommitteeMsgID), nil),
				},
				PostDutyRunnerStateRoot: "d9ccbe17862d695e9bceff0a9489c8fc81abc434d619a8efa97408b7e27ed4f9",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: failed decoding consensus data: incorrect offset",
			},
		},
	}
}
