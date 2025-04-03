package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// State holds all the relevant progress the duty execution progress
type State struct {
	PreConsensusContainer  *PartialSigContainer
	PostConsensusContainer *PartialSigContainer
	RunningInstance        *qbft.Instance
	DecidedValue           []byte
	// CurrentDuty is the duty the node pulled locally from the beacon node, might be different from decided duty
	StartingDuty types.Duty
	// flags
	Finished bool // Finished marked true when there is a full successful cycle (pre, consensus and post) with quorum
}

func NewRunnerState(quorum uint64, duty types.Duty) *State {
	return &State{
		PreConsensusContainer:  NewPartialSigContainer(quorum),
		PostConsensusContainer: NewPartialSigContainer(quorum),

		StartingDuty: duty,
		Finished:     false,
	}
}

// ReconstructBeaconSig aggregates collected partial beacon sigs
func (pcs *State) ReconstructBeaconSig(container *PartialSigContainer, root [32]byte, validatorPubKey []byte, validatorIndex phase0.ValidatorIndex) ([]byte, error) {
	// Reconstruct signatures
	signature, err := container.ReconstructSignature(root, validatorPubKey, validatorIndex)
	if err != nil {
		return nil, errors.Wrap(err, "could not reconstruct beacon sig")
	}
	return signature, nil
}
