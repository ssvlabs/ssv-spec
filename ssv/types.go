package ssv

import (
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/bloxapp/ssv-spec/p2p"
	"github.com/bloxapp/ssv-spec/types"
)

// DutyRunners is a map of duty runners mapped by msg id hex.
type DutyRunners map[types.BeaconRole]Runner

// DutyRunnerForMsgID returns a Runner from the provided msg ID, or nil if not found
func (ci DutyRunners) DutyRunnerForMsgID(msgID types.MessageID) Runner {
	role := msgID.GetRoleType()
	return ci[role]
}

// Network is the network interface for SSV
type Network interface {
	p2p.Broadcaster
}

// AttesterCalls interface has all attester duty specific calls
type AttesterCalls interface {
	// GetAttestationData returns attestation data by the given slot and committee index
	GetAttestationData(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (ssz.Marshaler, spec.DataVersion, error)
	// SubmitAttestation submit the attestation to the node
	SubmitAttestation(attestation *phase0.Attestation) error
}

// ProposerCalls interface has all block proposer duty specific calls
type ProposerCalls interface {
	// GetBeaconBlock returns beacon block by the given slot and committee index
	GetBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (ssz.Marshaler, spec.DataVersion, error)
	// GetBlindedBeaconBlock returns blinded beacon block by the given slot and committee index
	GetBlindedBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (ssz.Marshaler, spec.DataVersion, error)
	// SubmitBeaconBlock submit the block to the node
	SubmitBeaconBlock(block *spec.VersionedSignedBeaconBlock) error
	// SubmitBlindedBeaconBlock submit the blinded block to the node
	SubmitBlindedBeaconBlock(block *api.VersionedSignedBlindedBeaconBlock) error
}

// AggregatorCalls interface has all attestation aggregator duty specific calls
type AggregatorCalls interface {
	// SubmitAggregateSelectionProof returns an AggregateAndProof object
	SubmitAggregateSelectionProof(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeeLength uint64, index phase0.ValidatorIndex, slotSig []byte) (ssz.Marshaler, spec.DataVersion, error)
	// SubmitSignedAggregateSelectionProof broadcasts a signed aggregator msg
	SubmitSignedAggregateSelectionProof(msg *phase0.SignedAggregateAndProof) error
}

// SyncCommitteeCalls interface has all sync committee duty specific calls
type SyncCommitteeCalls interface {
	// GetSyncMessageBlockRoot returns beacon block root for sync committee
	GetSyncMessageBlockRoot(slot phase0.Slot) (phase0.Root, spec.DataVersion, error)
	// SubmitSyncMessage submits a signed sync committee msg
	SubmitSyncMessage(msg *altair.SyncCommitteeMessage) error
}

// SyncCommitteeContributionCalls interface has all sync committee contribution duty specific calls
type SyncCommitteeContributionCalls interface {
	// IsSyncCommitteeAggregator returns true if aggregator
	IsSyncCommitteeAggregator(proof []byte) (bool, error)
	// SyncCommitteeSubnetID returns sync committee subnet ID from subcommittee index
	SyncCommitteeSubnetID(index phase0.CommitteeIndex) (uint64, error)
	// GetSyncCommitteeContribution returns a types.Contributions object
	GetSyncCommitteeContribution(slot phase0.Slot, selectionProofs []phase0.BLSSignature, subnetIDs []uint64) (ssz.Marshaler, spec.DataVersion, error)
	// SubmitSignedContributionAndProof broadcasts to the network
	SubmitSignedContributionAndProof(contribution *altair.SignedContributionAndProof) error
}

type DomainCalls interface {
	DomainData(epoch phase0.Epoch, domain phase0.DomainType) (phase0.Domain, error)
}

type BeaconNode interface {
	// GetBeaconNetwork returns the beacon network the node is on
	GetBeaconNetwork() types.BeaconNetwork
	AttesterCalls
	ProposerCalls
	AggregatorCalls
	SyncCommitteeCalls
	SyncCommitteeContributionCalls
	DomainCalls
}
