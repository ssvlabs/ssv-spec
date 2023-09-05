package pre_consensus_justifications

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// JustificationsWithNoPrecon is a spec test that tests the scenario where message have ConsensusData with preconsensus justifications.
// But the tasks at hand require no preconsensus stage.
func JustificationsWithNoPrecon() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	consensusMsgs := func(cd *types.ConsensusData, role types.BeaconRole, height qbft.Height) []*types.SSVMessage {
		id := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], role)
		qbftMsgs := testingutils.SSVDecidingMsgsForHeight(cd, id[:], height, ks)

		ret := make([]*types.SSVMessage, 0)
		for _, msg := range qbftMsgs {
			byts, _ := msg.Encode()
			ret = append(ret, &types.SSVMessage{
				MsgType: types.SSVConsensusMsgType,
				MsgID:   id,
				Data:    byts,
			})
		}
		return ret
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus justifications for duties with no pre consensus",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: consensusMsgs(testingutils.TestSyncCommitteeConsensusDataWithPreconJust(ks),
					types.BNRoleSyncCommittee,
					testingutils.TestingDutySlot)[:1], //proposal message

				PostDutyRunnerStateRoot: "48c73f57659b69131467ef133ccb35d7de2fe96438d30bfa2b5ea63b19ead011",
				ExpectedError: "failed processing consensus message: could not process msg: invalid signed message" +
					": proposal not justified: proposal fullData invalid: invalid value: sync committee invalid justifications",
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: consensusMsgs(testingutils.TestAttesterConsensusDataWithPreconJust(ks), types.BNRoleAttester,
					testingutils.TestingDutySlot)[:1], //proposal message
				PostDutyRunnerStateRoot: "9d55ff5721b21c5b99dd4b4bacb0acda0b674112fe3cec55cc6aeb04ad5dc2fc",
				ExpectedError:           "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: invalid value: attester invalid justifications",
			},
		},
	}
}
