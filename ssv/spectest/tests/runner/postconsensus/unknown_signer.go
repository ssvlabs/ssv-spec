package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests an unknown  signer
func UnknownSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedError := "failed processing post consensus message: invalid post-consensus message: failed to verify PartialSignature: unknown signer"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus unknown signer",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: decideRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSigSyncCommitteeContributionWrongSignerMsg(ks.Shares[1], 5, 5, ks)),
				},
				PostDutyRunnerStateRoot: "5cfbb232547e630d3b293543d83730b4c6a3f28e51b42c86a730875de2903259",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "sync committee",
				Runner: decideRunner(
					testingutils.SyncCommitteeRunner(ks),
					&testingutils.TestingSyncCommitteeDuty,
					testingutils.TestSyncCommitteeConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSigSyncCommitteeWrongBeaconSignerMsg(ks.Shares[1], 5, 5)),
				},
				PostDutyRunnerStateRoot: "547f70ddf97fe30830513d655e3446ea5e6dd414498e16fa4c1039bd58d5bef9",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "aggregator",
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					&testingutils.TestingAggregatorDuty,
					testingutils.TestAggregatorConsensusData,
				),
				Duty: &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusSigAggregatorWrongBeaconSignerMsg(ks.Shares[1], 5, 5)),
				},
				PostDutyRunnerStateRoot: "55ab9040e690b441c218934d9a2ec3c2fae91107c020e769dc2f85130a59f595",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "attester",
				Runner: decideRunner(
					testingutils.AttesterRunner(ks),
					&testingutils.TestingAttesterDuty,
					testingutils.TestAttesterConsensusData,
				),
				Duty: &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusSigAttestationWrongBeaconSignerMsg(ks.Shares[1], 5, 5, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "11bcf29afa4d32c81018b38a0c63712b4f417101c6592b24d32ae6fbd7bdaded",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("proposer (%s)", version.String()),
			Runner: decideRunner(
				testingutils.ProposerRunner(ks),
				testingutils.TestingProposerDutyV(version),
				testingutils.TestProposerConsensusDataV(version),
			),
			Duty: testingutils.TestingProposerDutyV(version),
			Messages: []*types.SSVMessage{
				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusSigProposerWrongBeaconSignerMsgV(ks.Shares[1], 5, 5, version)),
			},
			PostDutyRunnerStateRoot: unknownSignerProposerSC(version).Root(),
			PostDutyRunnerState:     unknownSignerProposerSC(version).ExpectedState,
			OutputMessages:          []*types.SignedPartialSignatureMessage{},
			BeaconBroadcastedRoots:  []string{},
			DontStartDuty:           true,
			ExpectedError:           expectedError,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: decideRunner(
				testingutils.ProposerBlindedBlockRunner(ks),
				testingutils.TestingProposerDutyV(version),
				testingutils.TestProposerBlindedBlockConsensusDataV(version),
			),
			Duty: testingutils.TestingProposerDutyV(version),
			Messages: []*types.SSVMessage{
				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusSigProposerWrongBeaconSignerMsgV(ks.Shares[1], 5, 5, version)),
			},
			PostDutyRunnerStateRoot: unknownSignerBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     unknownSignerBlindedProposerSC(version).ExpectedState,
			OutputMessages:          []*types.SignedPartialSignatureMessage{},
			BeaconBroadcastedRoots:  []string{},
			DontStartDuty:           true,
			ExpectedError:           expectedError,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
