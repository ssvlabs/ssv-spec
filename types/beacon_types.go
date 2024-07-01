package types

import (
	"math"
	"time"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

var GenesisValidatorsRoot = spec.Root{}
var GenesisForkVersion = spec.Version{0, 0, 0, 0}

var (
	DomainProposer                    = [4]byte{0x00, 0x00, 0x00, 0x00}
	DomainAttester                    = [4]byte{0x01, 0x00, 0x00, 0x00}
	DomainRandao                      = [4]byte{0x02, 0x00, 0x00, 0x00}
	DomainDeposit                     = [4]byte{0x03, 0x00, 0x00, 0x00}
	DomainVoluntaryExit               = [4]byte{0x04, 0x00, 0x00, 0x00}
	DomainSelectionProof              = [4]byte{0x05, 0x00, 0x00, 0x00}
	DomainAggregateAndProof           = [4]byte{0x06, 0x00, 0x00, 0x00}
	DomainSyncCommittee               = [4]byte{0x07, 0x00, 0x00, 0x00}
	DomainSyncCommitteeSelectionProof = [4]byte{0x08, 0x00, 0x00, 0x00}
	DomainContributionAndProof        = [4]byte{0x09, 0x00, 0x00, 0x00}
	DomainApplicationBuilder          = [4]byte{0x00, 0x00, 0x00, 0x01}

	DomainError = [4]byte{0x99, 0x99, 0x99, 0x99}
)

// MaxEffectiveBalanceInGwei is the max effective balance
const MaxEffectiveBalanceInGwei uint64 = 32000000000

// BLSWithdrawalPrefixByte is the BLS withdrawal prefix
const BLSWithdrawalPrefixByte = byte(0)

// DefaultGasLimit sets gas limit used in validator registrations.
const DefaultGasLimit = 30_000_000

// BeaconRole type of the validator role for a specific duty
type BeaconRole uint64

const (
	BNRoleAttester BeaconRole = iota
	BNRoleAggregator
	BNRoleProposer
	BNRoleSyncCommittee
	BNRoleSyncCommitteeContribution

	BNRoleValidatorRegistration
	BNRoleVoluntaryExit

	BNRoleUnknown = math.MaxUint64
)

// String returns name of the role
func (r BeaconRole) String() string {
	switch r {
	case BNRoleAttester:
		return "ATTESTER"
	case BNRoleAggregator:
		return "AGGREGATOR"
	case BNRoleProposer:
		return "PROPOSER"
	case BNRoleSyncCommittee:
		return "SYNC_COMMITTEE"
	case BNRoleSyncCommitteeContribution:
		return "SYNC_COMMITTEE_CONTRIBUTION"
	case BNRoleValidatorRegistration:
		return "VALIDATOR_REGISTRATION"
	case BNRoleVoluntaryExit:
		return "VOLUNTARY_EXIT"
	default:
		return "UNDEFINED"
	}
}

type Duty interface {
	DutySlot() spec.Slot
	RunnerRole() RunnerRole
}

// ValidatorDuty represent data regarding the duty type with the duty data
type ValidatorDuty struct {
	// Type is the duty type (attest, propose)
	Type BeaconRole
	// PubKey is the public key of the validator that should attest.
	PubKey spec.BLSPubKey `ssz-size:"48"`
	// Slot is the slot in which the validator should attest.
	Slot spec.Slot
	// ValidatorIndex is the index of the validator that should attest.
	ValidatorIndex spec.ValidatorIndex
	// CommitteeIndex is the index of the committee in which the attesting validator has been placed.
	CommitteeIndex spec.CommitteeIndex
	// CommitteeLength is the length of the committee in which the attesting validator has been placed.
	CommitteeLength uint64
	// CommitteesAtSlot is the number of committees in the slot.
	CommitteesAtSlot uint64
	// ValidatorCommitteeIndex is the index of the validator in the list of validators in the committee.
	ValidatorCommitteeIndex uint64
	// ValidatorSyncCommitteeIndices is the index of the validator in the list of validators in the committee.
	ValidatorSyncCommitteeIndices []uint64 `ssz-max:"13"`
}

func MapDutyToRunnerRole(dutyRole BeaconRole) RunnerRole {
	switch dutyRole {
	case BNRoleAttester, BNRoleSyncCommittee:
		return RoleCommittee
	case BNRoleProposer:
		return RoleProposer
	case BNRoleAggregator:
		return RoleAggregator
	case BNRoleSyncCommitteeContribution:
		return RoleSyncCommitteeContribution
	case BNRoleValidatorRegistration:
		return RoleValidatorRegistration
	case BNRoleVoluntaryExit:
		return RoleVoluntaryExit
	}
	return RoleUnknown
}

func (bd *ValidatorDuty) DutySlot() spec.Slot {
	return bd.Slot
}

func (bd *ValidatorDuty) RunnerRole() RunnerRole {
	return MapDutyToRunnerRole(bd.Type)
}

// GetValidatorIndex returns the validator index
func (bd *ValidatorDuty) GetValidatorIndex() spec.ValidatorIndex {
	return bd.ValidatorIndex
}

type CommitteeDuty struct {
	Slot            spec.Slot
	ValidatorDuties []*ValidatorDuty
}

func (cd *CommitteeDuty) DutySlot() spec.Slot {
	return cd.Slot
}

func (cd *CommitteeDuty) RunnerRole() RunnerRole {
	return RoleCommittee
}

//

// Available networks.
const (
	// MainNetwork represents the main network.
	MainNetwork BeaconNetwork = "mainnet"

	// HoleskyNetwork represents the Holesky test network.
	HoleskyNetwork BeaconNetwork = "holesky"

	// PraterNetwork represents the Prater test network.
	PraterNetwork BeaconNetwork = "prater"

	// BeaconTestNetwork is a simple test network with a custom genesis time
	BeaconTestNetwork BeaconNetwork = "now_test_network"
)

// BeaconNetwork represents the network.
type BeaconNetwork string

// NetworkFromString returns network from the given string value
func NetworkFromString(n string) BeaconNetwork {
	switch n {
	case string(MainNetwork):
		return MainNetwork
	case string(HoleskyNetwork):
		return HoleskyNetwork
	case string(PraterNetwork):
		return PraterNetwork
	case string(BeaconTestNetwork):
		return BeaconTestNetwork
	default:
		return ""
	}
}

// ForkVersion returns the fork version of the network.
func (n BeaconNetwork) ForkVersion() [4]byte {
	switch n {
	case MainNetwork:
		return [4]byte{0, 0, 0, 0}
	case HoleskyNetwork:
		return [4]byte{0x01, 0x01, 0x70, 0x00}
	case PraterNetwork:
		return [4]byte{0x00, 0x00, 0x10, 0x20}
	case BeaconTestNetwork:
		return [4]byte{0x99, 0x99, 0x99, 0x99}
	default:
		return [4]byte{0x98, 0x98, 0x98, 0x98}
	}
}

// MinGenesisTime returns min genesis time value
func (n BeaconNetwork) MinGenesisTime() uint64 {
	switch n {
	case MainNetwork:
		return 1606824023
	case HoleskyNetwork:
		return 1695902400
	case PraterNetwork:
		return 1616508000
	case BeaconTestNetwork:
		return 1616508000
	default:
		return 0
	}

}

// SlotDurationSec returns slot duration
func (n BeaconNetwork) SlotDurationSec() time.Duration {
	return 12 * time.Second
}

// SlotsPerEpoch returns number of slots per one epoch
func (n BeaconNetwork) SlotsPerEpoch() uint64 {
	return 32
}

// EstimatedCurrentSlot returns the estimation of the current slot
func (n BeaconNetwork) EstimatedCurrentSlot() spec.Slot {
	return n.EstimatedSlotAtTime(time.Now().Unix())
}

// EstimatedSlotAtTime estimates slot at the given time
func (n BeaconNetwork) EstimatedSlotAtTime(time int64) spec.Slot {
	genesis := int64(n.MinGenesisTime())
	if time < genesis {
		return 0
	}
	return spec.Slot(uint64(time-genesis) / uint64(n.SlotDurationSec().Seconds()))
}

func (n BeaconNetwork) EstimatedTimeAtSlot(slot spec.Slot) int64 {
	d := int64(slot) * int64(n.SlotDurationSec().Seconds())
	return int64(n.MinGenesisTime()) + d
}

// EstimatedCurrentEpoch estimates the current epoch
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#compute_start_slot_at_epoch
func (n BeaconNetwork) EstimatedCurrentEpoch() spec.Epoch {
	return n.EstimatedEpochAtSlot(n.EstimatedCurrentSlot())
}

// EstimatedEpochAtSlot estimates epoch at the given slot
func (n BeaconNetwork) EstimatedEpochAtSlot(slot spec.Slot) spec.Epoch {
	return spec.Epoch(slot / spec.Slot(n.SlotsPerEpoch()))
}

func (n BeaconNetwork) FirstSlotAtEpoch(epoch spec.Epoch) spec.Slot {
	return spec.Slot(uint64(epoch) * n.SlotsPerEpoch())
}

func (n BeaconNetwork) EpochStartTime(epoch spec.Epoch) time.Time {
	firstSlot := n.FirstSlotAtEpoch(epoch)
	t := n.EstimatedTimeAtSlot(firstSlot)
	return time.Unix(t, 0)
}

// ComputeETHDomain returns computed domain
func ComputeETHDomain(domain spec.DomainType, fork spec.Version, genesisValidatorRoot spec.Root) (spec.Domain, error) {
	ret := spec.Domain{}
	copy(ret[0:4], domain[:])

	forkData := spec.ForkData{
		CurrentVersion:        fork,
		GenesisValidatorsRoot: genesisValidatorRoot,
	}
	forkDataRoot, err := forkData.HashTreeRoot()
	if err != nil {
		return ret, err
	}
	copy(ret[4:32], forkDataRoot[0:28])
	return ret, nil
}

func ComputeETHSigningRoot(obj ssz.HashRoot, domain spec.Domain) (spec.Root, error) {
	root, err := obj.HashTreeRoot()
	if err != nil {
		return spec.Root{}, err
	}
	signingContainer := spec.SigningData{
		ObjectRoot: root,
		Domain:     domain,
	}
	return signingContainer.HashTreeRoot()
}
