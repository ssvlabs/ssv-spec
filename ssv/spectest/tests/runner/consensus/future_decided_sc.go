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

// futureDecidedProposerSC returns state comparison object for the FutureDecided Proposer versioned spec test
func futureDecidedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
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
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version)),
					}),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				StartValue: cdBytes,
				State: &qbft.State{
					Share:  testingutils.TestingShare(testingutils.Testing4SharesSet()),
					ID:     ret.GetBaseRunner().QBFTController.Identifier,
					Round:  qbft.FirstRound,
					Height: qbft.FirstHeight,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})

			decidedInstance := &qbft.Instance{
				State: &qbft.State{
					Share:        testingutils.TestingShare(testingutils.Testing4SharesSet()),
					ID:           ret.GetBaseRunner().QBFTController.Identifier,
					Round:        1,
					Height:       2,
					Decided:      true,
					DecidedValue: testingutils.TestingQBFTFullData,
				},
			}
			comparable.SetMessages(
				decidedInstance,
				[]*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							ret.GetBaseRunner().QBFTController.Identifier,
						),
						nil,
					),
				},
			)

			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, decidedInstance)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = 2

			return ret
		}(),
	}
}

// futureDecidedBlindedProposerSC returns state comparison object for the FutureDecided BlindedProposer versioned spec test
func futureDecidedBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
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
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version)),
						testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version)),
					}),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				StartValue: cdBytes,
				State: &qbft.State{
					Share:  testingutils.TestingShare(testingutils.Testing4SharesSet()),
					ID:     ret.GetBaseRunner().QBFTController.Identifier,
					Round:  qbft.FirstRound,
					Height: qbft.FirstHeight,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})

			decidedInstance := &qbft.Instance{
				State: &qbft.State{
					Share:        testingutils.TestingShare(testingutils.Testing4SharesSet()),
					ID:           ret.GetBaseRunner().QBFTController.Identifier,
					Round:        1,
					Height:       2,
					Decided:      true,
					DecidedValue: testingutils.TestingQBFTFullData,
				},
			}
			comparable.SetMessages(
				decidedInstance,
				[]*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							ret.GetBaseRunner().QBFTController.Identifier,
						),
						nil,
					),
				},
			)

			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, decidedInstance)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = 2

			return ret
		}(),
	}
}
