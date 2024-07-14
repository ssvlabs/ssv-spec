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

// FutureDecided tests a running instance at a certain height, then processing a decided msg from a larger height.
// then returning an error and don't move to post consensus as it's not the same instance decided
func FutureDecided() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	getID := func(role types.RunnerRole) []byte {
		if role == types.RoleCommittee {
			opIDs := make([]types.OperatorID, len(ks.Committee()))
			for i, member := range ks.Committee() {
				opIDs[i] = member.Signer
			}
			committeeID := types.GetCommitteeID(opIDs)
			ret := types.NewMsgID(testingutils.TestingSSVDomainType, committeeID[:], role)
			return ret[:]
		}
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], role)
		return ret[:]
	}

	errStr := "failed processing consensus message: decided wrong instance"
	errStrCommittee := "no runner found for message's slot"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "consensus future decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "attester",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						testingutils.TestingDutySlot+1,
						getID(types.RoleCommittee),
					),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  errStrCommittee,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						testingutils.TestingDutySlot+1,
						getID(types.RoleCommittee),
					),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  errStrCommittee,
			},
			{
				Name:   "attester and sync committee",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterAndSyncCommitteeDuties,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						testingutils.TestingDutySlot+1,
						getID(types.RoleCommittee),
					),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  errStrCommittee,
			},
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						testingutils.TestingDutySlot+1,
						getID(types.RoleSyncCommitteeContribution),
					),
				},
				PostDutyRunnerStateRoot: futureDecidedSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     futureDecidedSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						testingutils.TestingDutySlot+1,
						getID(types.RoleAggregator),
					),
				},
				PostDutyRunnerStateRoot: futureDecidedAggregatorSC().Root(),
				PostDutyRunnerState:     futureDecidedAggregatorSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
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
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version))),
				testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
					[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
					[]types.OperatorID{1, 2, 3},
					qbft.Height(testingutils.TestingDutySlotV(version)+1),
					getID(types.RoleProposer),
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
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version))),
				testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
					[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
					[]types.OperatorID{1, 2, 3},
					qbft.Height(testingutils.TestingDutySlotV(version)+1),
					getID(types.RoleProposer),
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
