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

// quorum7OperatorsProposerSC returns state comparison object for the Quorum7Operators Proposer versioned spec test
func quorum7OperatorsProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing7SharesSet()
	cd := testingutils.TestProposerConsensusDataV(ks, version)
	cdBytes := testingutils.TestProposerConsensusDataBytsV(ks, version)

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				StartingDuty: &cd.Duty,
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], ks.Shares[4], 4, 4, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], ks.Shares[5], 5, 5, version)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				StartValue: cdBytes,
				State: &qbft.State{
					Share:  testingutils.TestingShare(ks),
					ID:     ret.GetBaseRunner().QBFTController.Identifier,
					Round:  qbft.FirstRound,
					Height: qbft.FirstHeight,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// quorum7OperatorsBlindedProposerSC returns state comparison object for the Quorum7Operators Blinded Proposer versioned spec test
func quorum7OperatorsBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing7SharesSet()
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(ks, version)
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(ks, version)

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				StartingDuty: &cd.Duty,
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], ks.Shares[4], 4, 4, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], ks.Shares[5], 5, 5, version)),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				StartValue: cdBytes,
				State: &qbft.State{
					Share:  testingutils.TestingShare(ks),
					ID:     ret.GetBaseRunner().QBFTController.Identifier,
					Round:  qbft.FirstRound,
					Height: qbft.FirstHeight,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}
