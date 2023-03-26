package postconsensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidMessageSlot tests a valid SignedPartialSignatureMessage with an invalid msg slot
func InvalidMessageSlot() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	invalidateSlot := func(msg *types.SignedPartialSignatureMessage) *types.SignedPartialSignatureMessage {
		msg.Message.Slot = testingutils.TestingDutySlot2
		return msg
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus invalid msg slot",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: decideRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, invalidateSlot(testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))),
				},
				PostDutyRunnerStateRoot: "f58387d4d4051a2de786e4cbf9dc370a8b19a544f52af04f71195feb3863fc5c",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing post consensus message: invalid post-consensus message: invalid partial sig slot",
			},
			{
				Name: "sync committee",
				Runner: decideRunner(
					testingutils.SyncCommitteeRunner(ks),
					&testingutils.TestingSyncCommitteeDuty,
					testingutils.TestSyncCommitteeConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(nil, invalidateSlot(testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: "599f535071e53121470fc10c80fad5d103340eba90dcd9672cff3e7a874de276",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing post consensus message: invalid post-consensus message: invalid partial sig slot",
			},
			{
				Name: "proposer",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					&testingutils.TestingProposerDuty,
					testingutils.TestProposerConsensusData,
				),
				Duty: &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, invalidateSlot(testingutils.PostConsensusProposerMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: "95dc99e33aadafa490275fc2d0361068a3f4cf7de901166606b2dbcc656b13c8",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing post consensus message: invalid post-consensus message: invalid partial sig slot",
			},
			{
				Name: "proposer (blinded block)",
				Runner: decideRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					&testingutils.TestingProposerDuty,
					testingutils.TestProposerBlindedBlockConsensusData,
				),
				Duty: &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, invalidateSlot(testingutils.PostConsensusProposerMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: "db66cd09c759130d5128e8a26e80dd3a16d3aa234097c8149d8b90c525eb15c9",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing post consensus message: invalid post-consensus message: invalid partial sig slot",
			},
			{
				Name: "aggregator",
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					&testingutils.TestingAggregatorDuty,
					testingutils.TestAggregatorConsensusData,
				),
				Duty: &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, invalidateSlot(testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: "1fb182fb19e446d61873abebc0ac85a3a9637b51d139cdbd7d8cb70cf7ffec82",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing post consensus message: invalid post-consensus message: invalid partial sig slot",
			},
			{
				Name: "attester",
				Runner: decideRunner(
					testingutils.AttesterRunner(ks),
					&testingutils.TestingAttesterDuty,
					testingutils.TestAttesterConsensusData,
				),
				Duty: &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, invalidateSlot(testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight))),
				},
				PostDutyRunnerStateRoot: "f43a47e0cb007d990f6972ce764ec8d0a35ae9c14a46f41bd7cde3df7d0e5f88",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing post consensus message: invalid post-consensus message: invalid partial sig slot",
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgValidatorRegistration(nil, invalidateSlot(testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight))),
				},
				PostDutyRunnerStateRoot: "f36c8b537afaba0894dbc8c87cb94466d8ac2623e9283f1c584e3d544b5f2b88",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no post consensus phase for validator registration",
			},
		},
	}
}
