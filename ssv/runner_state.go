package ssv

import (
	"crypto/sha256"
	"encoding/json"

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

// GetRoot returns the root used for signing and verification
func (pcs *State) GetRoot() ([32]byte, error) {
	marshaledRoot, err := pcs.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode State")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// Encode returns the encoded struct in bytes or error
func (pcs *State) Encode() ([]byte, error) {
	return json.Marshal(pcs)
}

// Decode returns error if decoding failed
func (pcs *State) Decode(data []byte) error {
	return json.Unmarshal(data, &pcs)
}

func (pcs *State) MarshalJSON() ([]byte, error) {

	// Create alias without duty
	type StateAlias struct {
		PreConsensusContainer  *PartialSigContainer
		PostConsensusContainer *PartialSigContainer
		RunningInstance        *qbft.Instance
		DecidedValue           []byte
		Finished               bool
		BeaconDuty             *types.BeaconDuty    `json:"BeaconDuty,omitempty"`
		CommitteeDuty          *types.CommitteeDuty `json:"CommitteeDuty,omitempty"`
	}

	alias := &StateAlias{
		PreConsensusContainer:  pcs.PreConsensusContainer,
		PostConsensusContainer: pcs.PostConsensusContainer,
		RunningInstance:        pcs.RunningInstance,
		DecidedValue:           pcs.DecidedValue,
		Finished:               pcs.Finished,
	}

	if pcs.StartingDuty != nil {
		if beaconDuty, ok := pcs.StartingDuty.(*types.BeaconDuty); ok {
			alias.BeaconDuty = beaconDuty
		} else if committeeDuty, ok := pcs.StartingDuty.(*types.CommitteeDuty); ok {
			alias.CommitteeDuty = committeeDuty
		} else {
			return nil, errors.New("can't marshal because BaseRunner.State.StartingDuty isn't BeaconDuty or CommitteeDuty")
		}
	}
	byts, err := json.Marshal(alias)

	return byts, err
}

func (pcs *State) UnmarshalJSON(data []byte) error {

	// Create alias without duty
	type StateAlias struct {
		PreConsensusContainer  *PartialSigContainer
		PostConsensusContainer *PartialSigContainer
		RunningInstance        *qbft.Instance
		DecidedValue           []byte
		Finished               bool
		BeaconDuty             *types.BeaconDuty    `json:"BeaconDuty,omitempty"`
		CommitteeDuty          *types.CommitteeDuty `json:"CommitteeDuty,omitempty"`
	}

	aux := &StateAlias{}

	// Unmarshal the JSON data into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	pcs.PreConsensusContainer = aux.PreConsensusContainer
	pcs.PostConsensusContainer = aux.PostConsensusContainer
	pcs.RunningInstance = aux.RunningInstance
	pcs.DecidedValue = aux.DecidedValue
	pcs.Finished = aux.Finished

	// Determine which type of duty was marshaled
	if aux.BeaconDuty != nil {
		pcs.StartingDuty = aux.BeaconDuty
	} else if aux.CommitteeDuty != nil {
		pcs.StartingDuty = aux.CommitteeDuty
	}

	return nil
}
