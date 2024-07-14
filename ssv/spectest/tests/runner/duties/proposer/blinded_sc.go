package proposer

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	ssvcomparable "github.com/ssvlabs/ssv-spec/ssv/spectest/comparable"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// fullHappyFlowProposerReceivingBlindedBlockSC returns state comparison object for the FullHappyFlow for a normal Proposer receiving a Blinded Block versioned spec test
func fullHappyFlowProposerReceivingBlindedBlockSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	// consensus data used for message creation
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
	// encoded consensus data relative to messages
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)
	// encoded consensus data that the runner set as StartValue
	startedCdBytes := testingutils.TestProposerConsensusDataBytsV(version)

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
				StartValue: startedCdBytes,
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

// fullHappyFlowBlindedProposerReceivingNormalBlockSC returns state comparison object for the FullHappyFlow for a Blinded Proposer receiving a Normal Block versioned spec test
func fullHappyFlowBlindedProposerReceivingNormalBlockSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	// consensus data used for message creation
	cd := testingutils.TestProposerConsensusDataV(version)
	// encoded consensus data relative to messages
	cdBytes := testingutils.TestProposerConsensusDataBytsV(version)
	// encoded consensus data that the runner set as StartValue
	startedCdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)

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
				StartValue: startedCdBytes,
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
