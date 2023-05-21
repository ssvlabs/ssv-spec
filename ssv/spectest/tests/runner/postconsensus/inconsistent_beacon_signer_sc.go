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

// inconsistentBeaconSignerSyncCommitteeContributionSC returns state comparison object for the InconsistentBeaconSigner SyncCommitteeContribution versioned spec test
func inconsistentBeaconSignerSyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeContributionConsensusData
	cdBytes := testingutils.TestSyncCommitteeContributionConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.SyncCommitteeContributionRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
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

// inconsistentBeaconSignerSyncCommitteeSC returns state comparison object for the InconsistentBeaconSigner SyncCommittee versioned spec test
func inconsistentBeaconSignerSyncCommitteeSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeConsensusData
	cdBytes := testingutils.TestSyncCommitteeConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.SyncCommitteeRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
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

// inconsistentBeaconSignerAggregatorSC returns state comparison object for the InconsistentBeaconSigner Aggregator versioned spec test
func inconsistentBeaconSignerAggregatorSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestAggregatorConsensusData
	cdBytes := testingutils.TestAggregatorConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.AggregatorRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
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

// inconsistentBeaconSignerAttesterSC returns state comparison object for the InconsistentBeaconSigner Attester versioned spec test
func inconsistentBeaconSignerAttesterSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestAttesterConsensusData
	cdBytes := testingutils.TestAttesterConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.AttesterRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
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

// inconsistentBeaconSignerProposerSC returns state comparison object for the InconsistentBeaconSigner Proposer versioned spec test
func inconsistentBeaconSignerProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerConsensusDataV(version)
	cdBytes := testingutils.TestProposerConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
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

// inconsistentBeaconSignerBlindedProposerSC returns state comparison object for the InconsistentBeaconSigner Blinded Proposer versioned spec test
func inconsistentBeaconSignerBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
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
