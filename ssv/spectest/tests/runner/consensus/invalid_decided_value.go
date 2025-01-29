package consensus

import (
	"crypto/rsa"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/qbft"

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

	expectedErr := "failed processing consensus message: decided ValidatorConsensusData invalid: decided value is invalid" +
		": invalid value: unknown duty role"
	expectedCommitteeErr := "failed processing consensus message: decided ValidatorConsensusData invalid: decided value" +
		" is invalid: attestation data source >= target"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "consensus decided invalid value",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(
						nil,
						testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					)),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(
						nil,
						testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2),
					)),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(
						nil,
						testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3),
					)),

					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.SyncCommitteeContributionMsgID,
						consensusDataByts(),
					),
				},
				PostDutyRunnerStateRoot: "aff4af0dbbead81d6cb9dd4ff734d4660712a8b5ab8e9016a3f86e2c2ead7549",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedErr,
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
				PostDutyRunnerStateRoot: "0c965c41a9318297ad03c27e79ca2d2d0fee357fff21995c014182ce5e2970b3",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedErr,
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
				PostDutyRunnerStateRoot: "d03cef76867bcb4540191a8e93a735b460ce5844271f718508a3821c404331a2",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedErr,
			},
		},
	}

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("aggregator (%s)", version.String()),
			Runner: testingutils.AggregatorRunner(ks),
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
					testingutils.AggregatorMsgID,
					consensusDataByts(),
				),
			},
			PostDutyRunnerStateRoot: "a93047858c5597f2b1de078a566e6b0227a217e10758a741ff7a2ed9e0a87d96",
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version),
			},
			ExpectedError: expectedErr,
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
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedCommitteeErr,
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
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedCommitteeErr,
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
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedCommitteeErr,
			},
		}...)
	}

	return multiSpecTest
}
