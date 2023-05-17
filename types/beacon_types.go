package types

import (
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

// BeaconRole type of the validator role for a specific duty
type BeaconRole uint64

// List of roles
const (
	BNRoleAttester BeaconRole = iota
	BNRoleAggregator
	BNRoleProposer
	BNRoleSyncCommittee
	BNRoleSyncCommitteeContribution

	BNRoleValidatorRegistration
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
	default:
		return "UNDEFINED"
	}
}

// Duty represent data regarding the duty type with the duty data
type Duty struct {
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
