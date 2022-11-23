package newduty

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// PostInvalidDecided tests starting a new duty after prev was decided with an invalid decided value
func PostInvalidDecided() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()

	consensusDataByts := func(role types.BeaconRole) []byte {
		cd := &types.ConsensusData{
			Duty: &types.Duty{
				Type:                    100,
				PubKey:                  testingutils.TestingValidatorPubKey,
				Slot:                    testingutils.TestingDutySlot,
				ValidatorIndex:          testingutils.TestingValidatorIndex,
				CommitteeIndex:          3,
				CommitteesAtSlot:        36,
				CommitteeLength:         128,
				ValidatorCommitteeIndex: 11,
			},
		}
		byts, _ := cd.Encode()
		return byts
	}

	decideWrong := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningHeight = qbft.FirstHeight
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
		err := r.GetBaseRunner().QBFTController.SaveInstance(qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight))
		if err != nil {
			panic(err)
		}

		err = r.ProcessConsensus(testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: r.GetBaseRunner().QBFTController.Identifier,
				Data:       testingutils.CommitDataBytes(consensusDataByts(r.GetBaseRunner().BeaconRoleType)),
			}))

		expectedError := "failed processing consensus message: decided ConsensusData invalid: decided value is invalid: duty invalid: wrong beacon role type"
		if err.Error() != expectedError {
			panic(err.Error())
		}
		return r
	}

	preInstances := func(r ssv.Runner, duty *types.Duty) []*qbft.Instance {
		runner := decideWrong(r, duty) // mock the process consensus in order to get the final version of the instance
		inst, err := runner.GetBaseRunner().QBFTController.InstanceForHeight(qbft.FirstHeight)
		if err != nil {
			panic(err.Error())
		}
		return []*qbft.Instance{
			inst,
		}
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post invalid decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  decideWrong(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionDuty),
				PreStoredInstances:      preInstances(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "2d3ee0dd2174000c00ba4f42b2df800e379533642e2c337a4f7cdb3f4f5893c8",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  decideWrong(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				PreStoredInstances:      preInstances(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "c9170c67f7541e9b4b6209608638724cb5122a1d94cb1327e3dfabfcc1b26fd0",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  decideWrong(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty),
				PreStoredInstances:      preInstances(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty),
				Duty:                    testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "dd57fb513dc1914ed961ae708f5625e85c4283ac1244e37fac67cc911083dfd7",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  decideWrong(testingutils.ProposerRunner(ks), testingutils.TestingProposerDuty),
				PreStoredInstances:      preInstances(testingutils.ProposerRunner(ks), testingutils.TestingProposerDuty),
				Duty:                    testingutils.TestingProposerDuty,
				PostDutyRunnerStateRoot: "18f07ee6a1aadcc6e627278e3335db6a90f947c9960a6dadba8ba446d87cf4d4",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  decideWrong(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				PreStoredInstances:      preInstances(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "524d235900dba12349af3e1fc103544f1eb066bb452c046869dc8985977f4e52",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
