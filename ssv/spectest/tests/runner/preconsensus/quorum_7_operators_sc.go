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

// quorum7OperatorsSyncCommitteeContributionSC returns state comparison object for the Quorum7Operators SyncCommitteeContribution versioned spec test
func quorum7OperatorsSyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing7SharesSet()
	cd := testingutils.TestSyncCommitteeContributionConsensusData

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.SyncCommitteeContributionRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{
						testingutils.SSVMsgSyncCommitteeContribution(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgSyncCommitteeContribution(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgSyncCommitteeContribution(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], 3)),
						testingutils.SSVMsgSyncCommitteeContribution(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[4], 4)),
						testingutils.SSVMsgSyncCommitteeContribution(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[5], 5)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{},
				),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot,
					LastPreparedRound: qbft.NoRound,
					Decided:           false,
				},
				StartValue: comparable.NoErrorEncoding(comparable.FixIssue178(cd, spec.DataVersionDeneb)),
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SignedSSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = testingutils.TestingDutySlot

			return ret
		}(),
	}
}

// quorum7OperatorsAggregatorSC returns state comparison object for the Quorum7Operators Aggregator versioned spec test
func quorum7OperatorsAggregatorSC() *comparable.StateComparison {
	ks := testingutils.Testing7SharesSet()
	cd := testingutils.TestAggregatorConsensusData
	cdBytes := testingutils.TestAggregatorConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.AggregatorRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{
						testingutils.SSVMsgAggregator(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgAggregator(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgAggregator(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], 3)),
						testingutils.SSVMsgAggregator(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], 4)),
						testingutils.SSVMsgAggregator(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[5], 5)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{},
				),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot,
					LastPreparedRound: qbft.NoRound,
					Decided:           false,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SignedSSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = testingutils.TestingDutySlot

			return ret
		}(),
	}
}

// quorum7OperatorsValidatorRegistrationSC returns state comparison object for the Quorum7Operators ValidatorRegistration versioned spec test
func quorum7OperatorsValidatorRegistrationSC() *comparable.StateComparison {
	ks := testingutils.Testing7SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ValidatorRegistrationRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{
						testingutils.SSVMsgValidatorRegistration(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgValidatorRegistration(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgValidatorRegistration(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3)),
						testingutils.SSVMsgValidatorRegistration(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[4], 4)),
						testingutils.SSVMsgValidatorRegistration(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[5], 5)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{},
				),
				StartingDuty: &testingutils.TestingValidatorRegistrationDuty,
				Finished:     true,
			}

			return ret
		}(),
	}
}

// quorum7OperatorsVoluntaryExitSC returns state comparison object for the Quorum7Operators VoluntaryExit versioned spec test
func quorum7OperatorsVoluntaryExitSC() *comparable.StateComparison {
	ks := testingutils.Testing7SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.VoluntaryExitRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{
						testingutils.SSVMsgVoluntaryExit(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgVoluntaryExit(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgVoluntaryExit(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[3], 3)),
						testingutils.SSVMsgVoluntaryExit(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[4], 4)),
						testingutils.SSVMsgVoluntaryExit(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[5], 5)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{},
				),
				StartingDuty: &testingutils.TestingVoluntaryExitDuty,
				Finished:     true,
			}

			return ret
		}(),
	}
}

// quorum7OperatorsProposerSC returns state comparison object for the Quorum7Operators Proposer versioned spec test
func quorum7OperatorsProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing7SharesSet()
	cd := testingutils.TestProposerConsensusDataV(version)
	cdBytes := testingutils.TestProposerConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				StartingDuty: &cd.Duty,
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{
						testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], 1, version)),
						testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], 2, version)),
						testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], 3, version)),
						testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, version)),
						testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], 5, version)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{},
				),
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				StartValue: cdBytes,
				State: &qbft.State{
					Share:  testingutils.TestingShare(ks),
					ID:     ret.GetBaseRunner().QBFTController.Identifier,
					Round:  qbft.FirstRound,
					Height: qbft.Height(testingutils.TestingDutySlotV(version)),
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SignedSSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version))

			return ret
		}(),
	}
}

// quorum7OperatorsBlindedProposerSC returns state comparison object for the Quorum7Operators Blinded Proposer versioned spec test
func quorum7OperatorsBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing7SharesSet()
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				StartingDuty: &cd.Duty,
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{
						testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], 1, version)),
						testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], 2, version)),
						testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], 3, version)),
						testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, version)),
						testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], 5, version)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(5),
					[]*types.SignedSSVMessage{},
				),
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				StartValue: cdBytes,
				State: &qbft.State{
					Share:  testingutils.TestingShare(ks),
					ID:     ret.GetBaseRunner().QBFTController.Identifier,
					Round:  qbft.FirstRound,
					Height: qbft.Height(testingutils.TestingDutySlotV(version)),
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SignedSSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version))

			return ret
		}(),
	}
}
