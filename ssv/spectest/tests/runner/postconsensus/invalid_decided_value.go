package postconsensus

import (
	"crypto/rsa"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidDecidedValue tests an invalid decided value
func InvalidDecidedValue() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	consensusDataByts := func(role types.BeaconRole) []byte {
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

	expectedErr := "failed processing post consensus message: invalid post-consensus message: no decided value"
	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus decided invalid value",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "attester",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.CommitteeMsgID(ks),
						testingutils.TestWrongBeaconVoteByts,
					),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, testingutils.TestingDutySlot))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedErr,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.CommitteeMsgID(ks),
						testingutils.TestWrongBeaconVoteByts,
					),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedErr,
			},
			{
				Name:   "attester and sync committee",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterAndSyncCommitteeDuties,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.CommitteeMsgID(ks),
						testingutils.TestWrongBeaconVoteByts,
					),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, testingutils.TestingDutySlot))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedErr,
			},
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),

					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.SyncCommitteeContributionMsgID,
						consensusDataByts(types.BNRoleSyncCommitteeContribution),
					),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))), // no qbft msg to mock the missing decided value
				},
				PostDutyRunnerStateRoot: "aff4af0dbbead81d6cb9dd4ff734d4660712a8b5ab8e9016a3f86e2c2ead7549",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),

					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{
							ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
						},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.AggregatorMsgID,
						consensusDataByts(types.BNRoleAggregator),
					),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: "a93047858c5597f2b1de078a566e6b0227a217e10758a741ff7a2ed9e0a87d96",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
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
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.ProposerMsgID,
						consensusDataByts(types.BNRoleProposer),
					),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				PostDutyRunnerStateRoot: "0c965c41a9318297ad03c27e79ca2d2d0fee357fff21995c014182ce5e2970b3",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(testingutils.Testing4SharesSet().Shares[1], 1, spec.DataVersionDeneb),
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
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.ProposerMsgID,
						consensusDataByts(types.BNRoleProposer),
					),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				PostDutyRunnerStateRoot: "d03cef76867bcb4540191a8e93a735b460ce5844271f718508a3821c404331a2",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(testingutils.Testing4SharesSet().Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedErr,
			},
		},
	}
}
