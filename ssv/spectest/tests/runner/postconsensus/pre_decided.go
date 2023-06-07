package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreDecided tests a valid SignedPartialSignatureMessage sent before decided reached, should error
func PreDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	err := "failed processing post consensus message: invalid post-consensus message: no decided value"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus before decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
				},
				PostDutyRunnerStateRoot: "79b26cb78d0e60ec0a20e796880a8b49549413c6a5c2d7be3c670d3215a7e686",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				BeaconBroadcastedRoots: []string{},
				ExpectedError:          err,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "5b941fd17da3e3dc11003611dedc70f6ecb4f647f2951a47ed11990bcd1f5afc",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				ExpectedError:           err,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "1d5d82620ea3ee1ba9dd64c60c1e09cccc2404881b4303348581fac5d2d2f647",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				BeaconBroadcastedRoots: []string{},
				ExpectedError:          err,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "1d42abde3ed6e27699960aa7476bb672a5c9f74f466896560f35597b56083853",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				ExpectedError:           err,
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
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, version)),

				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
			},
			PostDutyRunnerStateRoot: preDecidedProposerSC(version).Root(),
			PostDutyRunnerState:     preDecidedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{},
			ExpectedError:          err,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SSVMessage{
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, version)),

				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
			},
			PostDutyRunnerStateRoot: preDecidedBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     preDecidedBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{},
			ExpectedError:          err,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
