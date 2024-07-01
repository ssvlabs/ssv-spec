package newduty

import (
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostInvalidDecided tests starting a new duty after prev was decided with an invalid decided value
func PostInvalidDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	consensusDataByts := func() []byte {
		cd := &types.ValidatorConsensusData{
			Duty: types.ValidatorDuty{
				Type:                    100, // invalid
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

	// https://github.com/ssvlabs/ssv-spec/issues/285. We initialize the runner with an impossible decided value.
	// Maybe we should ensure that `ValidateDecided()` doesn't let the runner enter this state and delete the test?
	decideWrong := func(r ssv.Runner, duty types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.CommitteeMember,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.DutySlot()))
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.DutySlot())

		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		decidedValue := sha256.Sum256(consensusDataByts())
		r.GetBaseRunner().State.RunningInstance.State.DecidedValue = decidedValue[:]

		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post invalid decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name: "sync committee aggregator",
				Runner: decideWrong(testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "4112802181d740f78b68b0c67e4220e689af2ac1011a51d0f4c10c4df315fac5",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  decideWrong(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "b0ec12e65623dd1a95203d4a0a753bb6758f6bed3141a467bb7955530ae35ded",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  decideWrong(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionDeneb)),
				Duty:                    testingutils.TestingProposerDutyNextEpochV(spec.DataVersionDeneb),
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "c002484c2c25f5d97f625b5923484a062bdadb4eb21be9715dd9ae454883d890",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  decideWrong(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDutyNextEpoch,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "c002484c2c25f5d97f625b5923484a062bdadb4eb21be9715dd9ae454883d890",
				OutputMessages:          []*types.PartialSignatureMessages{},
			},
			{
				Name:                    "sync committee",
				Runner:                  decideWrong(testingutils.CommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDutyNextEpoch,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "c002484c2c25f5d97f625b5923484a062bdadb4eb21be9715dd9ae454883d890",
				OutputMessages:          []*types.PartialSignatureMessages{},
			},
			{
				Name:                    "attester and sync committee",
				Runner:                  decideWrong(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterAndSyncCommitteeDuties),
				Duty:                    testingutils.TestingAttesterAndSyncCommitteeDutiesNextEpoch,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "c002484c2c25f5d97f625b5923484a062bdadb4eb21be9715dd9ae454883d890",
				OutputMessages:          []*types.PartialSignatureMessages{},
			},
		},
	}
}
