package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	ssvcomparable "github.com/bloxapp/ssv-spec/ssv/spectest/comparable"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// quorumSyncCommitteeContributionSC returns state comparison object for the Quorum SyncCommitteeContribution versioned spec test
func quorumSyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeContributionConsensusData

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.SyncCommitteeContributionRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
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
					Decided:           false,
				},
				StartValue: comparable.NoErrorEncoding(comparable.FixIssue178(cd, spec.DataVersionBellatrix)),
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// quorumAggregatorSC returns state comparison object for the Quorum Aggregator versioned spec test
func quorumAggregatorSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestAggregatorConsensusData
	cdBytes := testingutils.TestAggregatorConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.AggregatorRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
						testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
						testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
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
					Decided:           false,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// quorumValidatorRegistrationSC returns state comparison object for the Quorum ValidatorRegistration versioned spec test
func quorumValidatorRegistrationSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ValidatorRegistrationRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				StartingDuty: &testingutils.TestingValidatorRegistrationDuty,
				Finished:     true,
			}

			return ret
		}(),
	}
}

// quorumProposerSC returns state comparison object for the Quorum Proposer versioned spec test
func quorumProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerConsensusDataV(version)
	cdBytes := testingutils.TestProposerConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, version)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
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
					Decided:           false,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// quorumBlindedProposerSC returns state comparison object for the Quorum Blinded Proposer versioned spec test
func quorumBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, version)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
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
					Decided:           false,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}
