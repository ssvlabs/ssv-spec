package postconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	ssvcomparable "github.com/bloxapp/ssv-spec/ssv/spectest/comparable"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// duplicateMsgQuorumSyncCommitteeContributionSC returns state comparison object for the duplicateMsgQuorum SyncCommitteeContribution versioned spec test
// it ignores the invalid duplicate message and doesn't insert it to the container
func duplicateMsgQuorumSyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeContributionConsensusData
	cdBytes := testingutils.TestSyncCommitteeContributionConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.SyncCommitteeContributionRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
					},
				),
				DecidedValue: cd,
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.NoRound,
					Decided:           true,
					DecidedValue:      cdBytes,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// duplicateMsgQuorumSyncCommitteeSC returns state comparison object for the duplicateMsgQuorum SyncCommittee versioned spec test
// it ignores the invalid duplicate message and doesn't insert it to the container
func duplicateMsgQuorumSyncCommitteeSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeConsensusData
	cdBytes := testingutils.TestSyncCommitteeConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.SyncCommitteeRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
					},
				),
				DecidedValue: cd,
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.NoRound,
					Decided:           true,
					DecidedValue:      cdBytes,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// duplicateMsgQuorumAggregatorSC returns state comparison object for the duplicateMsgQuorum Aggregator versioned spec test
// it ignores the invalid duplicate message and doesn't insert it to the container
func duplicateMsgQuorumAggregatorSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestAggregatorConsensusData
	cdBytes := testingutils.TestAggregatorConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.AggregatorRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
					},
				),
				DecidedValue: cd,
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.NoRound,
					Decided:           true,
					DecidedValue:      cdBytes,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// duplicateMsgQuorumAttesterSC returns state comparison object for the duplicateMsgQuorum Attester versioned spec test
// it ignores the invalid duplicate message and doesn't insert it to the container
func duplicateMsgQuorumAttesterSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestAttesterConsensusData
	cdBytes := testingutils.TestAttesterConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.AttesterRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
					},
				),
				DecidedValue: cd,
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.NoRound,
					Decided:           true,
					DecidedValue:      cdBytes,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// duplicateMsgQuorumProposerSC returns state comparison object for the duplicateMsgQuorum Proposer versioned spec test
// it ignores the invalid duplicate message and doesn't insert it to the container
func duplicateMsgQuorumProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerConsensusDataV(version)
	cdBytes := testingutils.TestProposerConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
					},
				),
				DecidedValue: cd,
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.NoRound,
					Decided:           true,
					DecidedValue:      cdBytes,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// duplicateMsgQuorumBlindedProposerSC returns state comparison object for the duplicateMsgQuorum Blinded Proposer versioned spec test
// it ignores the invalid duplicate message and doesn't insert it to the container
func duplicateMsgQuorumBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
					},
				),
				DecidedValue: cd,
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.NoRound,
					Decided:           true,
					DecidedValue:      cdBytes,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}
