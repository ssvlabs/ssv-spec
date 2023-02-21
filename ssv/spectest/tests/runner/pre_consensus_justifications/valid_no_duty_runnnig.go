package pre_consensus_justifications

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidNoRunningDuty tests a valid pre-consensus justification for a runner that has no running duty
func ValidNoRunningDuty() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	msgF := func(obj *types.ConsensusData) *qbft.SignedMessage {
		fullData, _ := obj.Encode()
		root, _ := obj.HashTreeRoot()
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     1,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.Shares[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus justification valid no running duty",
		Tests: []*tests.MsgProcessingSpecTest{
			//{
			//	Name:   "sync committee aggregator selection proof",
			//	Runner: testingutils.SyncCommitteeContributionRunner(ks),
			//	Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			//	},
			//	PostDutyRunnerStateRoot: "8d9edd36c3634e54d76985ddb4fa80f3427b47ab7dfab6053e7a396ab5ee494f",
			//	OutputMessages: []*types.SignedPartialSignatureMessage{
			//		testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
			//	},
			//},
			//{
			//	Name:   "aggregator selection proof",
			//	Runner: testingutils.AggregatorRunner(ks),
			//	Duty:   &testingutils.TestingAggregatorDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			//	},
			//	PostDutyRunnerStateRoot: "c5d864ca6a4ede7fe637846d080e0fe2cf1f4597c463cbf9a675bfbb78eacfc5",
			//	OutputMessages: []*types.SignedPartialSignatureMessage{
			//		testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
			//	},
			//},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(testingutils.TestProposerWithJustificationsConsensusData(ks)), nil),
				},
				PostDutyRunnerStateRoot: "5ba69a8b0fa59b3afbc64dc85edc32e45347169fb420378a801c6ed8486340fd",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			//{
			//	Name:   "randao (blinded block)",
			//	Runner: testingutils.ProposerBlindedBlockRunner(ks),
			//	Duty:   &testingutils.TestingProposerDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			//	},
			//	PostDutyRunnerStateRoot: "66967c4a461039e82dd60ca2ccd13ba82691bb43d5835a2b45394bfb4c0bc0ef",
			//	OutputMessages: []*types.SignedPartialSignatureMessage{
			//		testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
			//	},
			//},
			//{
			//	Name:   "attester",
			//	Runner: testingutils.AttesterRunner(ks),
			//	Duty:   &testingutils.TestingAttesterDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgAttester(nil, testingutils.PreConsensusFailedMsg(ks.Shares[1], 1)),
			//	},
			//	PostDutyRunnerStateRoot: "1d42abde3ed6e27699960aa7476bb672a5c9f74f466896560f35597b56083853",
			//	OutputMessages:          []*types.SignedPartialSignatureMessage{},
			//	ExpectedError:           "no pre consensus sigs required for attester role",
			//},
			//{
			//	Name:   "sync committee",
			//	Runner: testingutils.SyncCommitteeRunner(ks),
			//	Duty:   &testingutils.TestingSyncCommitteeDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgSyncCommittee(nil, testingutils.PreConsensusFailedMsg(ks.Shares[1], 1)),
			//	},
			//	PostDutyRunnerStateRoot: "26ea6d64e3660d677bc4fe7f02d951b7ddd5df142204f089cd8706b426b8a0d9",
			//	OutputMessages:          []*types.SignedPartialSignatureMessage{},
			//	ExpectedError:           "no pre consensus sigs required for sync committee role",
			//},
			//{
			//	Name:   "validator registration",
			//	Runner: testingutils.ValidatorRegistrationRunner(ks),
			//	Duty:   &testingutils.TestingValidatorRegistrationDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
			//	},
			//	PostDutyRunnerStateRoot: "6258dff05d5c0d040ce20933dd433073ac5badd1deb9f277097c0ce9bc92a57f",
			//	OutputMessages: []*types.SignedPartialSignatureMessage{
			//		testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
			//	},
			//},
		},
	}
}
