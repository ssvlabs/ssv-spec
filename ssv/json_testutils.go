package ssv

import (
	"crypto/sha256"
	"encoding/json"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// This file adds, as testing utils, the Encode, Decode and GetRoot methods
// so that structures follow the types.Encoder and types.Root interface

// State
func (pcs *State) Encode() ([]byte, error) {
	return json.Marshal(pcs)
}

func (pcs *State) Decode(data []byte) error {
	return json.Unmarshal(data, &pcs)
}

func (pcs *State) GetRoot() ([32]byte, error) {
	marshaledRoot, err := pcs.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode State")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

func (pcs *State) MarshalJSON() ([]byte, error) {

	// Create alias without duty
	type StateAlias struct {
		PreConsensusContainer  *PartialSigContainer
		PostConsensusContainer *PartialSigContainer
		RunningInstance        *qbft.Instance
		DecidedValue           []byte
		Finished               bool
		ValidatorDuty          *types.ValidatorDuty `json:"ValidatorDuty,omitempty"`
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
		if validatorDuty, ok := pcs.StartingDuty.(*types.ValidatorDuty); ok {
			alias.ValidatorDuty = validatorDuty
		} else if committeeDuty, ok := pcs.StartingDuty.(*types.CommitteeDuty); ok {
			alias.CommitteeDuty = committeeDuty
		} else {
			return nil, errors.New("can't marshal because BaseRunner.State.StartingDuty isn't ValidatorDuty or CommitteeDuty")
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
		ValidatorDuty          *types.ValidatorDuty `json:"ValidatorDuty,omitempty"`
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
	if aux.ValidatorDuty != nil {
		pcs.StartingDuty = aux.ValidatorDuty
	} else if aux.CommitteeDuty != nil {
		pcs.StartingDuty = aux.CommitteeDuty
	}

	return nil
}

// Committee
func (c *Committee) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Committee) Decode(data []byte) error {
	return json.Unmarshal(data, &c)
}

func (c *Committee) GetRoot() ([32]byte, error) {
	marshaledRoot, err := c.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode state")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

func (c *Committee) MarshalJSON() ([]byte, error) {

	type CommitteeAlias struct {
		Runners            map[spec.Slot]*CommitteeRunner
		CommitteeMember    types.CommitteeMember
		Share              map[spec.ValidatorIndex]*types.Share
		HighestDutySlotMap map[types.BeaconRole]map[spec.ValidatorIndex]spec.Slot
	}

	// Create object and marshal
	alias := &CommitteeAlias{
		Runners:         c.Runners,
		CommitteeMember: c.CommitteeMember,
		Share:           c.Share,
	}

	byts, err := json.Marshal(alias)

	return byts, err
}

func (c *Committee) UnmarshalJSON(data []byte) error {

	type CommitteeAlias struct {
		Runners            map[spec.Slot]*CommitteeRunner
		CommitteeMember    types.CommitteeMember
		Share              map[spec.ValidatorIndex]*types.Share
		HighestDutySlotMap map[types.BeaconRole]map[spec.ValidatorIndex]spec.Slot
	}

	// Unmarshal the JSON data into the auxiliary struct
	aux := &CommitteeAlias{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Assign fields
	c.Runners = aux.Runners
	c.CommitteeMember = aux.CommitteeMember
	c.Share = aux.Share

	return nil
}

// Runners

// ProposerRunner
func (r *ProposerRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ProposerRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

func (r *ProposerRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode ProposerRunner")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// CommitteeRunner
func (cr CommitteeRunner) Encode() ([]byte, error) {
	return json.Marshal(cr)
}

func (cr CommitteeRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &cr)
}

func (cr CommitteeRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := cr.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode CommitteeRunner")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// AggregatorRunner
func (r *AggregatorRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AggregatorRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

func (r *AggregatorRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode AggregatorRunner")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// SyncCommitteeAggregatorRunner
func (r *SyncCommitteeAggregatorRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *SyncCommitteeAggregatorRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

func (r *SyncCommitteeAggregatorRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode SyncCommitteeAggregatorRunner")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// ValidatorRegistrationRunner
func (r *ValidatorRegistrationRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ValidatorRegistrationRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

func (r *ValidatorRegistrationRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode ValidatorRegistrationRunner")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// VoluntaryExitRunner
func (r *VoluntaryExitRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *VoluntaryExitRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

func (r *VoluntaryExitRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode VoluntaryExitRunner")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}
