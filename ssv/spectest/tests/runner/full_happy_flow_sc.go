package runner

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	ssvcomparable "github.com/ssvlabs/ssv-spec/ssv/spectest/comparable"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// fullHappyFlowSyncCommitteeContributionSC returns state comparison object for the FullHappyFlow SyncCommitteeContribution spec test
func fullHappyFlowSyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeContributionConsensusData
	cdBytes := testingutils.TestSyncCommitteeContributionConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.SyncCommitteeContributionRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleSyncCommitteeContribution)[:3],
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks))),
					},
				),
				DecidedValue: testingutils.EncodeConsensusDataTest(comparable.FixIssue178(cd, spec.DataVersionPhase0)),
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:                 testingutils.TestingCommitteeMember(ks),
					ID:                              ret.GetBaseRunner().QBFTController.Identifier,
					Round:                           qbft.FirstRound,
					Height:                          qbft.Height(testingutils.TestingDutySlot),
					LastPreparedRound:               qbft.FirstRound,
					LastPreparedValue:               cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes, qbft.Height(testingutils.TestingDutySlot))),
					Decided:                         true,
					DecidedValue:                    cdBytes,
				},
				StartValue: comparable.NoErrorEncoding(comparable.FixIssue178(cd, spec.DataVersionBellatrix)),
			}
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlot)
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleSyncCommitteeContribution)[3:10],
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
		ExpectedState: func() ssv.Runner {
			ret := testingutils.AggregatorRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleAggregator)[:3],
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3))),
					},
				),
				DecidedValue: testingutils.EncodeConsensusDataTest(comparable.FixIssue178(cd, spec.DataVersionPhase0)),
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:                 testingutils.TestingCommitteeMember(ks),
					ID:                              ret.GetBaseRunner().QBFTController.Identifier,
					Round:                           qbft.FirstRound,
					Height:                          qbft.Height(testingutils.TestingDutySlot),
					LastPreparedRound:               qbft.FirstRound,
					LastPreparedValue:               cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes, qbft.Height(testingutils.TestingDutySlot))),
					Decided:                         true,
					DecidedValue:                    cdBytes,
				},
				StartValue: cdBytes,
			}
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlot)
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleAggregator)[3:10],
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
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleProposer)[:3],
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
					},
				),
				DecidedValue: testingutils.EncodeConsensusDataTest(comparable.FixIssue178(cd, version)),
				StartingDuty: &cd.Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.Height(testingutils.TestingDutySlotV(version)),
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes, qbft.Height(testingutils.TestingDutySlotV(version)))),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version))
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleProposer)[3:10],
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
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleProposer)[:3],
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
					},
				),
				DecidedValue: testingutils.EncodeConsensusDataTest(comparable.FixIssue178(cd, version)),
				StartingDuty: &testingutils.TestProposerConsensusDataV(version).Duty,
				Finished:     true,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.Height(testingutils.TestingDutySlotV(version)),
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes, qbft.Height(testingutils.TestingDutySlotV(version)))),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version))
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleProposer)[3:10],
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
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ValidatorRegistrationRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3))),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				StartingDuty: &testingutils.TestingValidatorRegistrationDuty,
				Finished:     true,
			}
			return ret
		}(),
	}
}

// fullHappyFlowVoluntaryExitSC returns state comparison object for the FullHappyFlow VoluntaryExit spec test
func fullHappyFlowVoluntaryExitSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.VoluntaryExitRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[2], 2))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[3], 3))),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				StartingDuty: &testingutils.TestingVoluntaryExitDuty,
				Finished:     true,
			}
			return ret
		}(),
	}
}
