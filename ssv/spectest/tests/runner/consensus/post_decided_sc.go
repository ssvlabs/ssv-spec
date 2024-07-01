package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	ssvcomparable "github.com/ssvlabs/ssv-spec/ssv/spectest/comparable"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// postDecidedSyncCommitteeContributionSC returns runner with decided instance and an extra commit message.
// There are pre-consensus messages in the container but no post-consensus messages.
func postDecidedSyncCommitteeContributionSC() *comparable.StateComparison {
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
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: testingutils.EncodeConsensusDataTest(comparable.FixIssue178(cd, spec.DataVersionPhase0)),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlot))),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: comparable.NoErrorEncoding(comparable.FixIssue178(cd, spec.DataVersionBellatrix)),
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleSyncCommitteeContribution)[3:10],
					testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.SyncCommitteeContributionMsgID, cdBytes),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = testingutils.TestingDutySlot
			return ret
		}(),
	}
}

// postDecidedAggregatorSC returns runner with decided instance and an extra commit message.
// There are pre-consensus messages in the container but no post-consensus messages.
func postDecidedAggregatorSC() *comparable.StateComparison {
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
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: testingutils.EncodeConsensusDataTest(comparable.FixIssue178(cd, spec.DataVersionPhase0)),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlot))),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleAggregator)[3:10],
					testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.AggregatorMsgID, cdBytes),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = testingutils.TestingDutySlot
			return ret
		}(),
	}
}

// postDecidedProposerSC returns runner with decided instance and an extra commit message.
// There are pre-consensus messages in the container but no post-consensus messages.
func postDecidedProposerSC(version spec.DataVersion) *comparable.StateComparison {
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
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: testingutils.EncodeConsensusDataTest(comparable.FixIssue178(cd, version)),
				StartingDuty: &cd.Duty,
				Finished:     false,
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
						ks.OperatorKeys[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlotV(version)))),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleProposer)[3:10],
					testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4], types.OperatorID(4),
						qbft.Height(testingutils.TestingDutySlotV(version)), testingutils.ProposerMsgID, cdBytes),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version))
			return ret
		}(),
	}
}

// postDecidedBlindedProposerSC  returns runner with decided instance and an extra commit message.
// // There are pre-consensus messages in the container but no post-consensus messages.
func postDecidedBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
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
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: testingutils.EncodeConsensusDataTest(comparable.FixIssue178(cd, version)),
				StartingDuty: &cd.Duty,
				Finished:     false,
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
						ks.OperatorKeys[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlotV(version)))),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.RoleProposer)[3:10],
					testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.OperatorKeys[4],
						types.OperatorID(4), qbft.Height(testingutils.TestingDutySlotV(version)), testingutils.ProposerMsgID,
						cdBytes),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version))
			return ret
		}(),
	}
}
