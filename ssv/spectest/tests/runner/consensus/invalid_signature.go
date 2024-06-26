package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/qbft"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidSignature tests a consensus message with an invalid signature
func InvalidSignature() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedError := "SignedSSVMessage has an invalid signature: crypto/rsa: verification error"
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus invalid signature",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),
					// Invalid Message
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[2], testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.SyncCommitteeContributionMsgID, testingutils.TestSyncCommitteeContributionConsensusDataByts,
							qbft.Height(testingutils.TestingDutySlot)), nil)),
				},
				PostDutyRunnerStateRoot: "3430b48cc4265a27d9f99e03355810b09c129b9e3c6cc83b4e7916777d595b2f",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					// Invalid Message
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[2], testingutils.SSVMsgSyncCommittee(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.SyncCommitteeMsgID, testingutils.TestSyncCommitteeConsensusDataByts,
							qbft.Height(testingutils.TestingDutySlot)), nil)),
				},
				PostDutyRunnerStateRoot: "339b34e6f00899baf8740299d3d453aea60c30cc0e28912a6ebb3770f59fc9b8",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),
					// Invalid Message
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[2], testingutils.SSVMsgAggregator(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.AggregatorMsgID, testingutils.TestAggregatorConsensusDataByts,
							qbft.Height(testingutils.TestingDutySlot)), nil)),
				},
				PostDutyRunnerStateRoot: "7134f3bfe0c675263254aadb1f73e452454418290f411b891090b2c76c5ae428",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, spec.DataVersionDeneb))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, spec.DataVersionDeneb))),
					// Invalid Message
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[2], testingutils.SSVMsgProposer(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ProposerMsgID, testingutils.TestProposerConsensusDataBytsV(spec.DataVersionDeneb),
							qbft.Height(testingutils.TestingDutySlotV(spec.DataVersionDeneb))), nil)),
				},
				PostDutyRunnerStateRoot: "01694261367dff43d4e85ebbfb1dd5d5081c78832ff693d77f300e0f8ffee071",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedError,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, spec.DataVersionDeneb))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, spec.DataVersionDeneb))),
					// Invalid Message
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[2], testingutils.SSVMsgProposer(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ProposerMsgID, testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
							qbft.Height(testingutils.TestingDutySlotV(spec.DataVersionDeneb))), nil)),
				},
				PostDutyRunnerStateRoot: "755134a5855f13a1e4a15b4b2a034a172c21f88a2af99e247d2cf4818aea30fe",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedError,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					// Invalid Message
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[2], testingutils.SSVMsgAttester(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.AttesterMsgID, testingutils.TestAttesterConsensusDataByts,
							qbft.Height(testingutils.TestingDutySlot)), nil)),
				},
				PostDutyRunnerStateRoot: "e062f2e50b8b308e83278c2f771c9473f4415bd8a64975bc5f29b61b29bd33fd",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3))),
					// Invalid Message
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[2], testingutils.SSVMsgValidatorRegistration(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ValidatorRegistrationMsgID, testingutils.TestAttesterConsensusDataByts,
							qbft.Height(testingutils.TestingDutySlot)), nil)),
				},
				PostDutyRunnerStateRoot: "ec573732e70b70808972c43acb5ead6443cff06ba30d8abb51e37ac82ffe0727",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingValidatorRegistration),
				},
				ExpectedError: expectedError,
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[2], 2))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[3], 3))),
					// Invalid Message
					testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[2], testingutils.SSVMsgVoluntaryExit(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.VoluntaryExitMsgID, testingutils.TestAttesterConsensusDataByts,
							qbft.Height(testingutils.TestingDutySlot)), nil)),
				},
				PostDutyRunnerStateRoot: "ec573732e70b70808972c43acb5ead6443cff06ba30d8abb51e37ac82ffe0727",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedVoluntaryExit(ks)),
				},
				ExpectedError: expectedError,
			},
		},
	}
}
