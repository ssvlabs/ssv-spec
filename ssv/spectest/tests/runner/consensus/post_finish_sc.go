package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	ssvcomparable "github.com/bloxapp/ssv-spec/ssv/spectest/comparable"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// postFinishProposerSC returns state comparison object for the PostFinish Proposer versioned spec test
func postFinishProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerConsensusDataV(version)
	cdBytes := testingutils.TestProposerConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				StartingDuty: &cd.Duty,
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[0:3]),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				DecidedValue: comparable.FixIssue178(cd, version),
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				StartValue: comparable.NoErrorEncoding(cd),
				State: &qbft.State{
					Share:  testingutils.TestingShare(testingutils.Testing4SharesSet()),
					ID:     ret.GetBaseRunner().QBFTController.Identifier,
					Round:  qbft.FirstRound,
					Height: qbft.FirstHeight,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1),
						ret.GetBaseRunner().QBFTController.Identifier,
						cdBytes,
					),
					LastPreparedRound: 1,
					LastPreparedValue: cdBytes,
					Decided:           true,
					DecidedValue:      cdBytes,
				},
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[3:10],
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.ProposerMsgID,
							testingutils.TestProposerConsensusDataBytsV(version),
						), nil),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			return ret
		}(),
	}
}

// postFinishBlindedProposerSC returns state comparison object for the PostFinish Blinded Proposer versioned spec test
func postFinishBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				StartingDuty: &cd.Duty,
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[0:3]),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SSVMessage{},
				),
				DecidedValue: comparable.FixIssue178(cd, version),
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				StartValue: comparable.NoErrorEncoding(cd),
				State: &qbft.State{
					Share:  testingutils.TestingShare(testingutils.Testing4SharesSet()),
					ID:     ret.GetBaseRunner().QBFTController.Identifier,
					Round:  qbft.FirstRound,
					Height: qbft.FirstHeight,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1),
						ret.GetBaseRunner().QBFTController.Identifier,
						cdBytes,
					),
					LastPreparedRound: 1,
					LastPreparedValue: cdBytes,
					Decided:           true,
					DecidedValue:      cdBytes,
				},
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[3:10],
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.ProposerMsgID,
							testingutils.TestProposerBlindedBlockConsensusDataBytsV(version),
						), nil,
					),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			return ret
		}(),
	}
}
