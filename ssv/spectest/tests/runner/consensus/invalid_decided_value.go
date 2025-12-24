package consensus

import (
	"crypto/rsa"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidDecidedValue tests an invalid decided value ValidatorConsensusData.Validate() != nil (unknown duty role)
func InvalidDecidedValue() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	consensusDataByts := func() []byte {
		cd := &types.ValidatorConsensusData{
			Duty: types.ValidatorDuty{
				Type:                    100,
				PubKey:                  testingutils.TestingValidatorPubKey,
				Slot:                    testingutils.TestingDutySlot,
				ValidatorIndex:          testingutils.TestingValidatorIndex,
				CommitteeIndex:          3,
				CommitteesAtSlot:        36,
				CommitteeLength:         128,
				ValidatorCommitteeIndex: 11,
			},
			Version: spec.DataVersionPhase0,
		}
		byts, _ := cd.Encode()
		return byts
	}
	aggregatorCommmitteeConsensusDataByts := func() []byte {
		cd := &types.AggregatorCommitteeConsensusData{}
		byts, _ := cd.Encode()
		return byts
	}

	expectedErrCode := types.QBFTValueInvalidErrorCode
	expectedCommitteeErrCode := types.AttestationSourceNotLessThanTargetErrorCode
	expectedAggCommitteeErrCode := types.AggCommConsensusDataNoValidatorErrorCode

	sccSlot := testingutils.TestingSyncCommitteeContributionDuty.Slot
	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"consensus decided invalid value",
		testdoc.ConsensusInvalidDecidedValueDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(
						nil,
						testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[1], ks.Shares[1], 1, 1, sccSlot),
					)),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(
						nil,
						testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[2], ks.Shares[2], 2, 2, sccSlot),
					)),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(
						nil,
						testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[3], ks.Shares[3], 3, 3, sccSlot),
					)),

					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(sccSlot),
						testingutils.TestingAggregatorCommitteeMsgID[:],
						aggregatorCommmitteeConsensusDataByts(),
					),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[1], ks.Shares[1], 1, 1, sccSlot),
				},
				ExpectedErrorCode: expectedAggCommitteeErrCode,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, spec.DataVersionDeneb))),

					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlotV(spec.DataVersionDeneb)),
						testingutils.ProposerMsgID,
						consensusDataByts(),
					),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedErrorCode: expectedErrCode,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, spec.DataVersionDeneb))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, spec.DataVersionDeneb))),

					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlotV(spec.DataVersionDeneb)),
						testingutils.ProposerMsgID,
						consensusDataByts(),
					),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedErrorCode: expectedErrCode,
			},
		},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("aggregator (%s)", version.String()),
			Runner: testingutils.AggregatorCommitteeRunner(ks),
			Duty:   testingutils.TestingAggregatorDuty(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3, version))),

				testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
					[]*rsa.PrivateKey{
						ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
					},
					[]types.OperatorID{1, 2, 3},
					qbft.Height(testingutils.TestingDutySlotV(version)),
					testingutils.TestingAggregatorCommitteeMsgID[:],
					aggregatorCommmitteeConsensusDataByts(),
				),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version),
			},
			ExpectedErrorCode: expectedAggCommitteeErrCode,
		},
		)
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:   fmt.Sprintf("attester (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlotV(version)),
						testingutils.CommitteeMsgID(ks),
						testingutils.TestWrongBeaconVoteByts,
					),
				},
				ExpectedErrorCode: expectedCommitteeErrCode,
			},
			{
				Name:   fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlotV(version)),
						testingutils.CommitteeMsgID(ks),
						testingutils.TestWrongBeaconVoteByts,
					),
				},
				ExpectedErrorCode: expectedCommitteeErrCode,
			},
			{
				Name:   fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlotV(version)),
						testingutils.CommitteeMsgID(ks),
						testingutils.TestWrongBeaconVoteByts,
					),
				},
				ExpectedErrorCode: expectedCommitteeErrCode,
			},
		}...)
	}

	return multiSpecTest
}
