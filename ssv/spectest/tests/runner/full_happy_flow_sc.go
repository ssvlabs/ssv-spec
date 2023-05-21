package runner

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	ssvcomparable "github.com/bloxapp/ssv-spec/ssv/spectest/comparable"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// fullHappyFlowSyncCommitteeContributionSC returns state comparison object for the FullHappyFlow SyncCommitteeContribution spec test
func fullHappyFlowSyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeContributionConsensusData
	cdBytes := testingutils.TestSyncCommitteeContributionConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.SyncCommitteeContributionRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleSyncCommitteeContribution)[:3],
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
					},
				),
				DecidedValue: comparable.FixIssue178(cd, spec.DataVersionPhase0),
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
					),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: comparable.NoErrorEncoding(comparable.FixIssue178(cd, spec.DataVersionBellatrix)),
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleSyncCommitteeContribution)[3:10],
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			return ret
		}(),
	}
}

// fullHappyFlowSyncCommitteeSC returns state comparison object for the FullHappyFlow SyncCommittee spec test
func fullHappyFlowSyncCommitteeSC() *comparable.StateComparison {
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
					[]*types.SSVMessage{
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
					},
				),
				DecidedValue: comparable.FixIssue178(cd, spec.DataVersionPhase0),
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
					),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleSyncCommittee)[:7],
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			return ret
		}(),
	}
}

// fullHappyFlowAggregatorSC returns state comparison object for the FullHappyFlow Aggregator spec test
func fullHappyFlowAggregatorSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestAggregatorConsensusData
	cdBytes := testingutils.TestAggregatorConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.AggregatorRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleAggregator)[:3],
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
					},
				),
				DecidedValue: comparable.FixIssue178(cd, spec.DataVersionPhase0),
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
					),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleAggregator)[3:10],
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			return ret
		}(),
	}
}

// fullHappyFlowProposerSC returns state comparison object for the FullHappyFlow Proposer versioned spec test
func fullHappyFlowProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerConsensusDataV(version)
	cdBytes := testingutils.TestProposerConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[:3],
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
					},
				),
				DecidedValue: comparable.FixIssue178(cd, version),
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
					),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[3:10],
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			return ret
		}(),
	}
}

// fullHappyFlowBlindedProposerSC returns state comparison object for the FullHappyFlow BlindedProposer versioned spec test
func fullHappyFlowBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[:3],
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
					},
				),
				DecidedValue: comparable.FixIssue178(cd, version),
				StartingDuty: &testingutils.TestProposerConsensusDataV(version).Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
					),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[3:10],
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			return ret
		}(),
	}
}

// fullHappyFlowAttesterSC returns state comparison object for the FullHappyFlow Attester spec test
func fullHappyFlowAttesterSC() *comparable.StateComparison {
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
					[]*types.SSVMessage{
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
					},
				),
				DecidedValue: comparable.FixIssue178(cd, spec.DataVersionPhase0),
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.FirstHeight,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
					),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleAttester)[:7],
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			return ret
		}(),
	}
}

// fullHappyFlowValidatorRegistrationSC returns state comparison object for the FullHappyFlow ValidatorRegistration spec test
func fullHappyFlowValidatorRegistrationSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
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
