package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PastSlot tests justification with highestDecidedDutySlot >= data.ValidatorDuty.Slot (not valid)
func PastSlot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msgF := func(obj *types.ValidatorConsensusData, id []byte) *types.SignedSSVMessage {

		// copy duty to isolate this change from all other spec tests
		copyDuty := types.ValidatorDuty{}
		byts, err := obj.Duty.MarshalSSZ()
		if err != nil {
			panic(err.Error())
		}
		if err := copyDuty.UnmarshalSSZ(byts); err != nil {
			panic(err.Error())
		}
		copyDuty.Slot = 5
		data := &types.ValidatorConsensusData{
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
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.RoleSyncCommitteeContribution),
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
					testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData, ks, types.RoleAggregator),
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
					testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb), ks, types.RoleProposer),
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
					testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb), ks, types.RoleProposer),
					msgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
				),
				PostDutyRunnerStateRoot: "2d8b29b5bf69b642910a9caa8a9662c2e71e5e13f19a5d0b52cc932eb29f3969",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedErr,
			},
		},
	}
}
