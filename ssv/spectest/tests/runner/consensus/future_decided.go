package consensus

import (
	"crypto/rsa"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureDecided tests a running instance at a certain height, then processing a decided msg from a larger height.
// then returning an error and don't move to post consensus as it's not the same instance decided
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
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgSyncCommitteeContribution(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommitteeContribution(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], 3)),
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
						[]types.OperatorID{1, 2, 3},
						testingutils.TestingDutySlot+1,
						getID(types.BNRoleSyncCommitteeContribution),
					),
				},
				PostDutyRunnerStateRoot: futureDecidedSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     futureDecidedSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
						[]types.OperatorID{1, 2, 3},
						testingutils.TestingDutySlot+1,
						getID(types.BNRoleSyncCommittee),
					),
				},
				PostDutyRunnerStateRoot: futureDecidedSyncCommitteeSC().Root(),
				PostDutyRunnerState:     futureDecidedSyncCommitteeSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
				ExpectedError:           errStr,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgAggregator(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgAggregator(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgAggregator(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], 3)),
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
						[]types.OperatorID{1, 2, 3},
						testingutils.TestingDutySlot+1,
						getID(types.BNRoleAggregator),
					),
				},
				PostDutyRunnerStateRoot: futureDecidedAggregatorSC().Root(),
				PostDutyRunnerState:     futureDecidedAggregatorSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
						[]types.OperatorID{1, 2, 3},
						testingutils.TestingDutySlot+1,
						getID(types.BNRoleAttester),
					),
				},
				PostDutyRunnerStateRoot: futureDecidedAttesterSC().Root(),
				PostDutyRunnerState:     futureDecidedAttesterSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
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
			Messages: []*types.SignedSSVMessage{
				testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version)),
				testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
					[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
					[]types.OperatorID{1, 2, 3},
					qbft.Height(testingutils.TestingDutySlotV(version)+1),
					getID(types.BNRoleProposer),
				),
			},
			PostDutyRunnerStateRoot: futureDecidedProposerSC(version).Root(),
			PostDutyRunnerState:     futureDecidedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
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
			Messages: []*types.SignedSSVMessage{
				testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version)),
				testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
					[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
					[]types.OperatorID{1, 2, 3},
					qbft.Height(testingutils.TestingDutySlotV(version)+1),
					getID(types.BNRoleProposer),
				),
			},
			PostDutyRunnerStateRoot: futureDecidedBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     futureDecidedBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
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
