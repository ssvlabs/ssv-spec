package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FutureMessage tests a valid proposal future msg
func FutureMessage() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	futureMsgF := func(value types.Encoder, id []byte) *types.SignedSSVMessage {
		var fullData []byte
		if value != nil {
			fullData, _ = value.Encode()
		} else {
			panic("no consensus data or beacon vote")
		}
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     10,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.OperatorKeys[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	expectedError := "failed processing consensus message: future msg from height, could not process"
	expectedErrorCommittee := "no runner found for message's slot"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus future message",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "attester",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					futureMsgF(&testingutils.TestBeaconVote, testingutils.CommitteeMsgID(ks)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				DontStartDuty:  true,
				ExpectedError:  expectedErrorCommittee,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					futureMsgF(&testingutils.TestBeaconVote, testingutils.CommitteeMsgID(ks)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				DontStartDuty:  true,
				ExpectedError:  expectedErrorCommittee,
			},
			{
				Name:   "attester sync committee",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterAndSyncCommitteeDuties,
				Messages: []*types.SignedSSVMessage{
					futureMsgF(&testingutils.TestBeaconVote, testingutils.CommitteeMsgID(ks)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				DontStartDuty:  true,
				ExpectedError:  expectedErrorCommittee,
			},
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestSyncCommitteeContributionConsensusData, testingutils.SyncCommitteeContributionMsgID),
				},
				PostDutyRunnerStateRoot: "68fd25b1cb30902e7b7b3e7ff674c3862ff956954a06fac0df485961b8bb3934",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestAggregatorConsensusData, testingutils.AggregatorMsgID),
				},
				PostDutyRunnerStateRoot: "bdc7c2150e0f2d4669e112848f5140b52aba0367b60ff2b594d5a5bef3587834",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb), testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "32dd1d1d7a4c34bb7dafc0866f69eb49f6a0a23755b135f83ad14d12e39fff82",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
						testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "58b946451dc5ccbd52fbc9e6bbe0ac888253d1708be018a3ff0b07762dd28891",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1),
						testingutils.ValidatorRegistrationMsgID, testingutils.TestAttesterConsensusDataByts,
						qbft.Height(testingutils.TestingDutySlot)),
				},
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for validator registration",
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1),
						testingutils.VoluntaryExitMsgID, testingutils.TestAttesterConsensusDataByts,
						qbft.Height(testingutils.TestingDutySlot)),
				},
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for voluntary exit",
			},
		},
	}
}
