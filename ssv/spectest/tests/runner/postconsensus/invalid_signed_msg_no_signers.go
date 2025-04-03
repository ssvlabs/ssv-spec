package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidSignedMessageNoSigners tests an invalid SignedSSVMessage with no signers
func InvalidSignedMessageNoSigners() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	noSigners := func(msg *types.SignedSSVMessage) *types.SignedSSVMessage {
		msg.OperatorIDs = make([]uint64, 0)
		return msg
	}

	expectedError := "invalid SignedSSVMessage: no signers"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus invalid signed message no signers",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: decideRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					noSigners(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 0, ks)))),
				},
				PostDutyRunnerStateRoot: "f58387d4d4051a2de786e4cbf9dc370a8b19a544f52af04f71195feb3863fc5c",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "proposer",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					noSigners(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 0, spec.DataVersionDeneb)))),
				},
				PostDutyRunnerStateRoot: "ff213af6f0bf2350bb37f48021c137dd5552b1c25cb5c6ebd0c1d27debf6080e",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name: "proposer (blinded block)",
				Runner: decideRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					noSigners(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 0, spec.DataVersionDeneb)))),
				},
				PostDutyRunnerStateRoot: "9b4524d5100835df4d71d0a1e559acdc33d541c44a746ebda115c5e7f3eaa85a",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
		},
	}

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("aggregator (%s)", version.String()),
			Runner: decideRunner(
				testingutils.AggregatorRunner(ks),
				testingutils.TestingAggregatorDuty(version),
				testingutils.TestAggregatorConsensusData(version),
			),
			Duty: testingutils.TestingAggregatorDuty(version),
			Messages: []*types.SignedSSVMessage{
				noSigners(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 0, version)))),
			},
			PostDutyRunnerStateRoot: "1fb182fb19e446d61873abebc0ac85a3a9637b51d139cdbd7d8cb70cf7ffec82",
			OutputMessages:          []*types.PartialSignatureMessages{},
			BeaconBroadcastedRoots:  []string{},
			DontStartDuty:           true,
			ExpectedError:           expectedError,
		},
		)
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{

			{
				Name: fmt.Sprintf("attester (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterDuty(version),
				Messages: []*types.SignedSSVMessage{
					noSigners(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 0, version)))),
				},
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{},
				DontStartDuty:          true,
				ExpectedError:          expectedError,
			},
			{
				Name: fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingSyncCommitteeDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingSyncCommitteeDuty(version),
				Messages: []*types.SignedSSVMessage{
					noSigners(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 0, version)))),
				},
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{},
				DontStartDuty:          true,
				ExpectedError:          expectedError,
			},
			{
				Name: fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterAndSyncCommitteeDuties(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: []*types.SignedSSVMessage{
					noSigners(testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 0, version)))),
				},
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{},
				DontStartDuty:          true,
				ExpectedError:          expectedError,
			},
		}...)
	}

	return multiSpecTest
}
