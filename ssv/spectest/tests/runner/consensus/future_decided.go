package consensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureDecided tests a running instance at FirstHeight, then processing a decided msg from height 2 and returning decided but doesn't move to post consensus as it's not the same instance decided
func FutureDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	getID := func(role types.BeaconRole) []byte {
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], role)
		return ret[:]
	}

	errStr := "failed processing consensus message: decided wrong instance"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "consensus future decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleSyncCommitteeContribution),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "6e003cd912d95b73c177f276929f03180b02dc77fec38720d718f8361b53d2b8",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleSyncCommittee),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "ef4080c54e4d1cad83e02388ce61ce499c93628a5176d5defb12c97856dcb4ee",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           errStr,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					testingutils.SSVMsgAggregator(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleAggregator),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "96cb700dbfff1b478bb550f6b76e66e9837b05e78b77692ca9acc4bebb1845d7",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleAttester),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "fdc2d5c52ee8e5dc1a33219579aa820b2e8a22d75bf159cbd68f6299af546e26",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           errStr,
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SSVMessage{
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version)),
				testingutils.SSVMsgProposer(
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						2,
						getID(types.BNRoleProposer),
					),
					nil,
				),
			},
			PostDutyRunnerStateRoot: futureDecidedProposerSC(version).Root(),
			PostDutyRunnerState:     futureDecidedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			},
			ExpectedError: errStr,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SSVMessage{
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version)),
				testingutils.SSVMsgProposer(
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						2,
						getID(types.BNRoleProposer),
					),
					nil,
				),
			},
			PostDutyRunnerStateRoot: futureDecidedBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     futureDecidedBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			},
			ExpectedError: errStr,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
