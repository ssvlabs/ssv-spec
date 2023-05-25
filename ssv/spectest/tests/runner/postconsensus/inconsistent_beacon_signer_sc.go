package postconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	ssvcomparable "github.com/bloxapp/ssv-spec/ssv/spectest/comparable"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// inconsistentBeaconSignerProposerSC returns state comparison object for the InconsistentBeaconSigner Proposer versioned spec test
func inconsistentBeaconSignerProposerSC(version spec.DataVersion) *comparable.StateComparison {
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
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
				DecidedValue: cd,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:        testingutils.TestingShare(ks),
					ID:           ret.GetBaseRunner().QBFTController.Identifier,
					Round:        qbft.FirstRound,
					Height:       qbft.FirstHeight,
					Decided:      true,
					DecidedValue: cdBytes,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}

// inconsistentBeaconSignerBlindedProposerSC returns state comparison object for the InconsistentBeaconSigner Blinded Proposer versioned spec test
func inconsistentBeaconSignerBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
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
					[]*types.SSVMessage{},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSignatureContainer(),
					[]*types.SSVMessage{},
				),
				DecidedValue: cd,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:        testingutils.TestingShare(ks),
					ID:           ret.GetBaseRunner().QBFTController.Identifier,
					Round:        qbft.FirstRound,
					Height:       qbft.FirstHeight,
					Decided:      true,
					DecidedValue: cdBytes,
				},
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SSVMessage{})
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)

			return ret
		}(),
	}
}
