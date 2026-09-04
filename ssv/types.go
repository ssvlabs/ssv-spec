package ssv

import (
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/p2p"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// DutyRunners is a map of duty runners mapped by msg id hex.
type DutyRunners map[types.RunnerRole]Runner

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
	GetAttestationData(slot phase0.Slot) (*phase0.AttestationData,
		spec.DataVersion, error)
	// SubmitAttestation submit the attestation to the node
	SubmitAttestations(attestations []*spec.VersionedAttestation) error
}

// ProposerCalls interface has all block proposer duty specific calls
type ProposerCalls interface {
	// GetBeaconBlock returns beacon block by the given slot, graffiti, and randao.
	GetBeaconBlock(slot phase0.Slot, graffiti, randao []byte) (*api.VersionedProposal, ssz.Marshaler, error)
	// SubmitBeaconBlock submit the block to the node
	SubmitBeaconBlock(block *api.VersionedProposal, sig phase0.BLSSignature) error
	// GetGloasBeaconBlock returns the Gloas (ePBS) beacon block for the given slot, graffiti, and
	// randao (SIP #94 §4). Separate from GetBeaconBlock because api.VersionedProposal cannot carry a
	// Gloas block; there is no blinded variant — Gloas blocks are bid-only.
	GetGloasBeaconBlock(slot phase0.Slot, graffiti, randao []byte) (*gloas.BeaconBlock, error)
	// SubmitGloasBeaconBlock submits the signed Gloas (ePBS) block to the node (SIP #94 §4)
	SubmitGloasBeaconBlock(block *gloas.BeaconBlock, sig phase0.BLSSignature) error
}

// AggregatorCalls interface has all attestation aggregator duty specific calls
type AggregatorCalls interface {
	// IsAggregator returns true if the validator is selected as an aggregator
	IsAggregator(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeeLength uint64, slotSig []byte) bool
	// GetAggregateAttestation returns the aggregate attestation for the given slot and committee
	GetAggregateAttestation(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (ssz.Marshaler, error)
	// SubmitAggregateSelectionProof returns an AggregateAndProof object
	// Deprecated: Use IsAggregator and GetAggregateAttestation instead. Kept for backward compatibility.
	SubmitAggregateSelectionProof(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeeLength uint64, index phase0.ValidatorIndex, slotSig []byte) (ssz.Marshaler, spec.DataVersion, error)
	// SubmitSignedAggregateAndProof broadcasts a signed aggregator msg
	SubmitSignedAggregateAndProof(msg *spec.VersionedSignedAggregateAndProof) error
	// SubmitMultipleSignedAggregateAndProof broadcasts multiple signed aggregator msgs
	SubmitMultipleSignedAggregateAndProof(msg []*spec.VersionedSignedAggregateAndProof) error
}

// SyncCommitteeCalls interface has all sync committee duty specific calls
type SyncCommitteeCalls interface {
	// GetSyncMessageBlockRoot returns beacon block root for sync committee
	GetSyncMessageBlockRoot(slot phase0.Slot) (phase0.Root, spec.DataVersion, error)
	// SubmitSyncMessages submits a signed sync committee msg
	SubmitSyncMessages(msgs []*altair.SyncCommitteeMessage) error
}

// SyncCommitteeContributionCalls interface has all sync committee contribution duty specific calls
type SyncCommitteeContributionCalls interface {
	// IsSyncCommitteeAggregator returns true if aggregator
	IsSyncCommitteeAggregator(proof []byte) bool
	// SyncCommitteeSubnetID returns sync committee subnet ID from subcommittee index
	SyncCommitteeSubnetID(index phase0.CommitteeIndex) uint64
	// GetSyncCommitteeContribution returns a types.Contributions object
	GetSyncCommitteeContribution(slot phase0.Slot, selectionProofs []phase0.BLSSignature, subnetIDs []uint64) (ssz.Marshaler, spec.DataVersion, error)
	// SubmitSignedContributionAndProof broadcasts to the network
	SubmitSignedContributionAndProof(contribution *altair.SignedContributionAndProof) error
}

// ValidatorRegistrationCalls interface has all validator registration duty specific calls
type ValidatorRegistrationCalls interface {
	// SubmitValidatorRegistration submits a validator registration
	SubmitValidatorRegistration(registration *api.VersionedSignedValidatorRegistration) error
}

// VoluntaryExitCalls interface has all validator voluntary exit duty specific calls
type VoluntaryExitCalls interface {
	// SubmitVoluntaryExit submits a validator voluntary exit
	SubmitVoluntaryExit(voluntaryExit *phase0.SignedVoluntaryExit) error
}

type DomainCalls interface {
	DomainData(epoch phase0.Epoch, domain phase0.DomainType) (phase0.Domain, error)
}

type VersionCalls interface {
	// DataVersion returns a data version for the given epoch.
	// In practice, for performance, responses can be cached in order not to always trigger an API call.
	DataVersion(epoch phase0.Epoch) spec.DataVersion
}

// ProposerPreferencesCalls interface has all Gloas (ePBS) proposer-preferences duty specific calls (SIP #94 §5)
type ProposerPreferencesCalls interface {
	// ProposerDutiesDependentRoot returns the dependent root the epoch's proposer duties were computed
	// against; the preference pins it so consumers can tell which duty assignment it was emitted for.
	ProposerDutiesDependentRoot(epoch phase0.Epoch) (phase0.Root, error)
	// SubmitProposerPreferences submits the reconstructed signed proposer preferences to the node
	SubmitProposerPreferences(preferences *gloas.SignedProposerPreferences) error
}

// EnvelopeCalls interface has all Gloas (ePBS) execution-payload envelope duty specific calls (SIP #94 §6)
type EnvelopeCalls interface {
	// GetBlindedExecutionPayloadEnvelope returns this operator's own produced blinded envelope for the
	// slot's decided block. It answers only on the beacon node that built the block (held from the §4
	// produceBlockV4 self-build response), so an error means this operator is not the builder operator.
	GetBlindedExecutionPayloadEnvelope(slot phase0.Slot, blockRoot phase0.Root) (*gloas.BlindedExecutionPayloadEnvelope, error)
	// SubmitExecutionPayloadEnvelope publishes the reconstructed reveal. The builder operator passes the
	// blinded envelope plus the threshold-reconstructed signature (valid for the full envelope by
	// root-equivalence); its beacon node forms the full SignedExecutionPayloadEnvelope body from its
	// cache (SIP #94 §6).
	SubmitExecutionPayloadEnvelope(envelope *gloas.BlindedExecutionPayloadEnvelope, signature phase0.BLSSignature) error
}

// PTCCalls interface has all Gloas (ePBS) Payload Timeliness Committee duty specific calls (SIP #94 §3)
type PTCCalls interface {
	// GetPayloadAttestationData returns the slot's payload attestation data as observed by the local
	// beacon node; a zero BeaconBlockRoot signals no block was seen for the slot — the abstain trigger.
	GetPayloadAttestationData(slot phase0.Slot) (*gloas.PayloadAttestationData, error)
	// SubmitPayloadAttestation submits the reconstructed payload attestation message to the node
	SubmitPayloadAttestation(msg *gloas.PayloadAttestationMessage) error
}

type BeaconNode interface {
	// GetBeaconNetwork returns the beacon network the node is on
	GetBeaconNetwork() types.BeaconNetwork
	AttesterCalls
	ProposerCalls
	AggregatorCalls
	SyncCommitteeCalls
	SyncCommitteeContributionCalls
	ValidatorRegistrationCalls
	VoluntaryExitCalls
	PTCCalls
	ProposerPreferencesCalls
	EnvelopeCalls
	DomainCalls
	VersionCalls
}
