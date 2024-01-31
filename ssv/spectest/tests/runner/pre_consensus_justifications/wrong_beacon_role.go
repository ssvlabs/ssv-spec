package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongBeaconRole tests wrong ConsensusData.Duty.Type != runner.BeaconRoleType
// (changed the consensus data object between runner types, e.g. sync committee receives contribution proof)
func WrongBeaconRole() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msgF := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
		fullData, _ := obj.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     1,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.Shares[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	expectedErr := "failed processing consensus message: invalid pre-consensus justification: wrong beacon role"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus wrong beacon role",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: decideFirstHeight(testingutils.SyncCommitteeContributionRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(msgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID), nil),
				},
				PostDutyRunnerStateRoot: "2619aeecde47fe0efc36aa98fbb2df9834d9eee77f62abe0d10532dbd5215790",
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
					testingutils.SSVMsgAggregator(msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID), nil),
				},
				PostDutyRunnerStateRoot: "db1b416873d19be76cddc92ded0d442ba0e642514973b5dfec45f587c6ffde15",
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
					testingutils.SSVMsgProposer(msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "2754fc7ced14fb15f3f18556bb6b837620287cbbfbf908abafa5a0533fc4bc5f",
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
					testingutils.SSVMsgProposer(msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "6bd59da9f817b8e40112e58231e36738b9d021db4416c9eeec1dd0236a5362e2",
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
					testingutils.SSVMsgAttester(msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.AttesterMsgID), nil),
				},
				PostDutyRunnerStateRoot: "81cb7b1d3ea3087d49f9773b3a2b75a87b901e50427d237f2a10c0e1904e7684",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: duty invalid: wrong beacon role type",
			},
			{
				Name:   "sync committee",
				Runner: decideFirstHeight(testingutils.SyncCommitteeRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeMsgID), nil),
				},
				PostDutyRunnerStateRoot: "38592232077cd45709a7c6cfdd20c9d899af9d79bc750add3c4b8f2b6794cb34",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: duty invalid: wrong beacon role type",
			},
		},
	}
}
