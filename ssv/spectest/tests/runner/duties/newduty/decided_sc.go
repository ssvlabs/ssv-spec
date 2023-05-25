package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	ssvcomparable "github.com/bloxapp/ssv-spec/ssv/spectest/comparable"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// postDecidedProposerSC returns state comparison object for the PostDecided Proposer versioned spec test
func postDecidedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				StartingDuty: testingutils.TestingProposerDutyNextEpochV(version),
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
			}
			instance := &qbft.Instance{
				State: &qbft.State{
					Share:   testingutils.TestingShare(ks),
					ID:      ret.GetBaseRunner().QBFTController.Identifier,
					Round:   qbft.FirstRound,
					Height:  qbft.FirstHeight,
					Decided: true,
				},
			}
			comparable.SetMessages(instance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, instance)
			return ret
		}(),
	}
}

// postDecidedBlindedProposerSC returns state comparison object for the PostDecided Blinded Proposer versioned spec test
func postDecidedBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				StartingDuty: testingutils.TestingProposerDutyNextEpochV(version),
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
			}
			instance := &qbft.Instance{
				State: &qbft.State{
					Share:   testingutils.TestingShare(ks),
					ID:      ret.GetBaseRunner().QBFTController.Identifier,
					Round:   qbft.FirstRound,
					Height:  qbft.FirstHeight,
					Decided: true,
				},
			}
			comparable.SetMessages(instance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, instance)
			return ret
		}(),
	}
}
