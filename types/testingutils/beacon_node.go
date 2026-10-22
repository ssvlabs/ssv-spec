package testingutils

import (
	"encoding/hex"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// Testing Beacon Node
// ==================================================

type TestingBeaconNode struct {
	BroadcastedRoots             []phase0.Root
	SyncCommitteeAggregatorRoots map[string]bool
	CommitteeIndexAggregators    map[phase0.CommitteeIndex]bool
	// producesExternalBid makes GetBeaconBlock return a bare external-bid block and no envelope (SIP #94 §4).
	producesExternalBid bool
	// wrongPayloadAttestationSlot makes GetPayloadAttestationData answer for another slot (SIP #94 §3).
	wrongPayloadAttestationSlot bool
	// failBlockSubmit makes SubmitBeaconBlock fail without recording the block (SIP #94 §6 reveal-on-attempt).
	failBlockSubmit bool
}

func NewTestingBeaconNode() *TestingBeaconNode {
	return &TestingBeaconNode{
		BroadcastedRoots:             []phase0.Root{},
		SyncCommitteeAggregatorRoots: make(map[string]bool),
		CommitteeIndexAggregators:    make(map[phase0.CommitteeIndex]bool),
	}
}

// SetSyncCommitteeAggregatorRootHexes FOR TESTING ONLY!! sets which sync committee aggregator roots will return true for aggregator
func (bn *TestingBeaconNode) SetSyncCommitteeAggregatorRootHexes(roots map[string]bool) {
	bn.SyncCommitteeAggregatorRoots = roots
}

// SetAggregators FOR TESTING ONLY!! sets committee indices values for IsAggregator
func (bn *TestingBeaconNode) SetAggregators(committeeIndices map[phase0.CommitteeIndex]bool) {
	bn.CommitteeIndexAggregators = committeeIndices
}

// SetProducesExternalBid FOR TESTING ONLY!! sets producesExternalBid
func (bn *TestingBeaconNode) SetProducesExternalBid(v bool) {
	bn.producesExternalBid = v
}

// SetWrongPayloadAttestationSlot FOR TESTING ONLY!! sets wrongPayloadAttestationSlot
func (bn *TestingBeaconNode) SetWrongPayloadAttestationSlot(v bool) {
	bn.wrongPayloadAttestationSlot = v
}

// SetFailBlockSubmit FOR TESTING ONLY!! sets failBlockSubmit
func (bn *TestingBeaconNode) SetFailBlockSubmit(v bool) {
	bn.failBlockSubmit = v
}

// GetBeaconNetwork returns the beacon network the node is on
func (bn *TestingBeaconNode) GetBeaconNetwork() types.BeaconNetwork {
	return types.BeaconTestNetwork
}

// GetAttestationData returns attestation data by the given slot and committee index
func (bn *TestingBeaconNode) GetAttestationData(slot phase0.Slot) (*phase0.AttestationData, spec.DataVersion, error) {
	data := *TestingAttestationData(gloas.DataVersionGloas)
	data.Slot = slot
	return &data, gloas.DataVersionGloas, nil
}

// SubmitAttestations submit attestations to the node
// Note: The test is concerned with what should be sent on the wire. Thus, attestations are converted into a
// SingleAttestation object as in the Ethereum spec.
func (bn *TestingBeaconNode) SubmitAttestations(attestations []*spec.VersionedAttestation) error {
	for _, att := range attestations {
		if att.Version != gloas.DataVersionGloas {
			panic("unsupported version")
		}
		singleAttestation, err := att.Gloas.ToSingleAttestation(att.ValidatorIndex)
		if err != nil {
			panic(err)
		}
		root, _ := singleAttestation.HashTreeRoot()
		bn.BroadcastedRoots = append(bn.BroadcastedRoots, root)
	}
	return nil
}

// SubmitVoluntaryExit submit the VoluntaryExit object to the node
func (bn *TestingBeaconNode) SubmitVoluntaryExit(voluntaryExit *phase0.SignedVoluntaryExit) error {
	r, _ := voluntaryExit.HashTreeRoot()
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// GetBeaconBlock returns the self-build fixture block for the slot plus the blinded envelope derived from it, so
// this operator is the builder operator; under SetProducesExternalBid, a bare external-bid block and no envelope
// (SIP #94 §4/§6).
func (bn *TestingBeaconNode) GetBeaconBlock(slot phase0.Slot, graffiti, randao []byte) (*gloas.BeaconBlock, *gloas.BlindedExecutionPayloadEnvelope, error) {
	if bn.producesExternalBid {
		// External bid win: a bare block and no envelope, so the decided value carries a zero payload_root.
		return gloas.TestingBeaconBlockExternalBuild(slot), nil, nil
	}
	return gloas.TestingBeaconBlock(slot), TestingBlindedExecutionPayloadEnvelope(slot), nil
}

// SubmitBeaconBlock records the signed block's root.
func (bn *TestingBeaconNode) SubmitBeaconBlock(block *gloas.BeaconBlock, sig phase0.BLSSignature) error {
	if bn.failBlockSubmit {
		return fmt.Errorf("forced block submit failure")
	}
	sb := &gloas.SignedBeaconBlock{
		Message:   block,
		Signature: sig,
	}
	r, err := sb.HashTreeRoot()
	if err != nil {
		return err
	}
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// SubmitExecutionPayloadEnvelope records the published §6 reveal's root — the blinded envelope's root,
// which by root-equivalence is the full envelope's (SIP #94 §6).
func (bn *TestingBeaconNode) SubmitExecutionPayloadEnvelope(envelope *gloas.BlindedExecutionPayloadEnvelope, signature phase0.BLSSignature) error {
	r, err := envelope.HashTreeRoot()
	if err != nil {
		return err
	}
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// IsAggregator returns true if the validator is selected as an aggregator
func (bn *TestingBeaconNode) IsAggregator(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeeLength uint64, slotSig []byte) bool {
	// In production, this would check the selection proof against the committee modulo

	// Check if committee index is set
	if val, found := bn.CommitteeIndexAggregators[committeeIndex]; found {
		return val
	}

	// Always return true for testing, for committees not set
	return true
}

// GetAggregateAttestation returns the aggregate attestation for the given slot and committee
func (bn *TestingBeaconNode) GetAggregateAttestation(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (ssz.Marshaler, error) {
	return TestingGloasAggregateAndProofV(TestingValidatorIndex, gloas.DataVersionGloas).Aggregate, nil
}

// SubmitAggregateSelectionProof returns an AggregateAndProof object
// Deprecated: Use IsAggregator and GetAggregateAttestation instead. Kept for backward compatibility.
func (bn *TestingBeaconNode) SubmitAggregateSelectionProof(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeeLength uint64, index phase0.ValidatorIndex, slotSig []byte) (ssz.Marshaler, spec.DataVersion, error) {
	return TestingAggregateAndProofV(gloas.DataVersionGloas, TestingValidatorIndex), gloas.DataVersionGloas, nil
}

// SubmitSignedAggregateAndProof broadcasts a signed aggregator msg
func (bn *TestingBeaconNode) SubmitSignedAggregateAndProof(msg *spec.VersionedSignedAggregateAndProof) error {
	if msg.Version != gloas.DataVersionGloas {
		panic("unsupported version")
	}
	root, _ := msg.Gloas.HashTreeRoot()
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, root)
	return nil
}

// SubmitMultipleSignedAggregateAndProof broadcasts multiple signed aggregator msgs
func (bn *TestingBeaconNode) SubmitMultipleSignedAggregateAndProof(msg []*spec.VersionedSignedAggregateAndProof) error {
	for _, m := range msg {
		if err := bn.SubmitSignedAggregateAndProof(m); err != nil {
			return err
		}
	}
	return nil
}

// GetSyncMessageBlockRoot returns beacon block root for sync committee
func (bn *TestingBeaconNode) GetSyncMessageBlockRoot(slot phase0.Slot) (phase0.Root, spec.DataVersion, error) {
	return TestingSyncCommitteeBlockRoot, gloas.DataVersionGloas, nil
}

// SubmitSyncMessage submits a signed sync committee msg
func (bn *TestingBeaconNode) SubmitSyncMessages(msgs []*altair.SyncCommitteeMessage) error {
	for _, msg := range msgs {
		r, _ := msg.HashTreeRoot()
		bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	}
	return nil
}

// IsSyncCommitteeAggregator returns tru if aggregator
func (bn *TestingBeaconNode) IsSyncCommitteeAggregator(proof []byte) bool {
	if len(bn.SyncCommitteeAggregatorRoots) != 0 {
		if val, found := bn.SyncCommitteeAggregatorRoots[hex.EncodeToString(proof)]; found {
			return val
		}
		return false
	}
	return true
}

// SyncCommitteeSubnetID returns sync committee subnet ID from subcommittee index
func (bn *TestingBeaconNode) SyncCommitteeSubnetID(index phase0.CommitteeIndex) uint64 {
	// Real calculation:
	// Each subnet has syncCommitteeSize / subnetCount validators
	// subnetCount is 4 for mainnet
	// const subnetCount = 4
	// const syncCommitteeSize = 512
	// const subnetSize = syncCommitteeSize / subnetCount
	// Outputs index / subnetSize
	return uint64(index) / (512 / 4)
}

// GetSyncCommitteeContribution returns
func (bn *TestingBeaconNode) GetSyncCommitteeContribution(slot phase0.Slot, selectionProofs []phase0.BLSSignature, subnetIDs []uint64) (ssz.Marshaler, spec.DataVersion, error) {
	return &TestingContributionsData, gloas.DataVersionGloas, nil
}

// SubmitSignedContributionAndProof broadcasts to the network
func (bn *TestingBeaconNode) SubmitSignedContributionAndProof(contribution *altair.SignedContributionAndProof) error {
	r, _ := contribution.HashTreeRoot()
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

func (bn *TestingBeaconNode) DomainData(epoch phase0.Epoch, domain phase0.DomainType) (phase0.Domain, error) {
	// epoch is used to calculate fork version, here we hard code it
	return types.ComputeETHDomain(domain, types.GenesisForkVersion, types.GenesisValidatorsRoot)
}
