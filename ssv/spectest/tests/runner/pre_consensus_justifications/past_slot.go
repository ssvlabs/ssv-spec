package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastSlot tests justification with highestDecidedDutySlot >= data.BeaconDuty.Slot (not valid)
func PastSlot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msgF := func(obj *types.ConsensusData, id []byte) *types.SignedSSVMessage {

		// copy duty to isolate this change from all other spec tests
		copyDuty := types.BeaconDuty{}
		byts, err := obj.Duty.MarshalSSZ()
		if err != nil {
			panic(err.Error())
		}
		if err := copyDuty.UnmarshalSSZ(byts); err != nil {
			panic(err.Error())
		}
		copyDuty.Slot = 5
		data := &types.ConsensusData{
			Duty:                       copyDuty,
			Version:                    obj.Version,
			PreConsensusJustifications: obj.PreConsensusJustifications,
			DataSSZ:                    obj.DataSSZ,
		}

		fullData, _ := data.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     1,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.OperatorKeys[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	expectedErr := "failed processing consensus message: invalid pre-consensus justification: duty.slot <= highest decided slot"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus past slot",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution),
					msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID),
				),
				PostDutyRunnerStateRoot: "bcaaa6fda6ec6e09f355371b391499ef38c4cf0e19f2ee00560906a8cd65fc94",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator),
					msgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID),
				),
				PostDutyRunnerStateRoot: "21ec7a48e5132e13de7499e7d44e4d7665287106433fd32dfba016b1a1afd02e",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb), ks, types.BNRoleProposer),
					msgF(testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
				),
				PostDutyRunnerStateRoot: "62b2dd10ed329c8ffeb4ba4a4fd3597f9746d8cd03ea8a7907ac0e6c5638b1b3",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb), ks, types.BNRoleProposer),
					msgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
				),
				PostDutyRunnerStateRoot: "2d8b29b5bf69b642910a9caa8a9662c2e71e5e13f19a5d0b52cc932eb29f3969",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedErr,
			},
			{

				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester),
					msgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID),
				),
				PostDutyRunnerStateRoot: "16c1db7c756f2e7dfff270c3ce0f9b1ee321b28bc57b9394cf030d99d29f25a1",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee),
					msgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID),
				),
				PostDutyRunnerStateRoot: "618756a8dbfef101c4a18ffc128de49c3bcac17388575bd3ab58db9ec3ce1b71",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
		},
	}
}
