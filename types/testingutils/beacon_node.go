package testingutils

import (
	"encoding/hex"
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"
	apiv1fulu "github.com/attestantio/go-eth2-client/api/v1/fulu"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Testing Beacon Node
// ==================================================

type TestingBeaconNode struct {
	BroadcastedRoots             []phase0.Root
	syncCommitteeAggregatorRoots map[string]bool
}

func NewTestingBeaconNode() *TestingBeaconNode {
	return &TestingBeaconNode{
		BroadcastedRoots: []phase0.Root{},
	}
}

// SetSyncCommitteeAggregatorRootHexes FOR TESTING ONLY!! sets which sync committee aggregator roots will return true for aggregator
func (bn *TestingBeaconNode) SetSyncCommitteeAggregatorRootHexes(roots map[string]bool) {
	bn.syncCommitteeAggregatorRoots = roots
}

// GetBeaconNetwork returns the beacon network the node is on
func (bn *TestingBeaconNode) GetBeaconNetwork() types.BeaconNetwork {
	return types.BeaconTestNetwork
}

// GetAttestationData returns attestation data by the given slot and committee index
func (bn *TestingBeaconNode) GetAttestationData(slot phase0.Slot) (*phase0.AttestationData, spec.DataVersion, error) {
	version := VersionBySlot(slot)
	data := *TestingAttestationData(version)
	data.Slot = slot
	return &data, version, nil
}

// SubmitAttestations submit attestations to the node
// Note: The test is concerned with what should be sent on the wire. Thus, electra Attestations are converted into a SingleAttestation object as in the Ethereum spec.
func (bn *TestingBeaconNode) SubmitAttestations(attestations []*spec.VersionedAttestation) error {
	for _, att := range attestations {
		var root [32]byte

		switch att.Version {
		case spec.DataVersionPhase0:
			root, _ = att.Phase0.HashTreeRoot()
		case spec.DataVersionAltair:
			root, _ = att.Altair.HashTreeRoot()
		case spec.DataVersionBellatrix:
			root, _ = att.Bellatrix.HashTreeRoot()
		case spec.DataVersionCapella:
			root, _ = att.Capella.HashTreeRoot()
		case spec.DataVersionDeneb:
			root, _ = att.Deneb.HashTreeRoot()
		case spec.DataVersionElectra:
			singleAttestation, err := att.Electra.ToSingleAttestation(att.ValidatorIndex)
			if err != nil {
				panic(err)
			}
			root, _ = singleAttestation.HashTreeRoot()
		case spec.DataVersionFulu:
			singleAttestation, err := att.Fulu.ToSingleAttestation(att.ValidatorIndex)
			if err != nil {
				panic(err)
			}
			root, _ = singleAttestation.HashTreeRoot()
		default:
			panic("unsupported version")
		}

		bn.BroadcastedRoots = append(bn.BroadcastedRoots, root)
	}
	return nil
}

func (bn *TestingBeaconNode) SubmitValidatorRegistration(registration *api.VersionedSignedValidatorRegistration) error {
	r, _ := registration.V1.HashTreeRoot()
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// SubmitVoluntaryExit submit the VoluntaryExit object to the node
func (bn *TestingBeaconNode) SubmitVoluntaryExit(voluntaryExit *phase0.SignedVoluntaryExit) error {
	r, _ := voluntaryExit.HashTreeRoot()
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// GetBeaconBlock returns beacon block by the given slot, graffiti, and randao.
func (bn *TestingBeaconNode) GetBeaconBlock(slot phase0.Slot, graffiti, randao []byte) (*api.VersionedProposal, ssz.Marshaler, error) {
	version := VersionBySlot(slot)
	vBlk := TestingBeaconBlockV(version)

	switch version {
	case spec.DataVersionCapella:
		return vBlk, vBlk.Capella, nil
	case spec.DataVersionDeneb:
		return vBlk, vBlk.Deneb, nil
	case spec.DataVersionElectra:
		return vBlk, vBlk.Electra, nil
	case spec.DataVersionFulu:
		return vBlk, vBlk.Fulu, nil
	default:
		return nil, nil, fmt.Errorf("unsupported version %s", version)
	}
}

// SubmitBeaconBlock submit the block (blinded or full) to the node
func (bn *TestingBeaconNode) SubmitBeaconBlock(block *api.VersionedProposal, sig phase0.BLSSignature) error {
	var r [32]byte

	if block.Blinded {
		// Handle blinded blocks
		switch block.Version {
		case spec.DataVersionCapella:
			if block.CapellaBlinded == nil {
				return fmt.Errorf("%s blinded block is nil", block.Version.String())
			}
			sb := &apiv1capella.SignedBlindedBeaconBlock{
				Message:   block.CapellaBlinded,
				Signature: sig,
			}
			r, _ = sb.HashTreeRoot()
		case spec.DataVersionDeneb:
			if block.DenebBlinded == nil {
				return fmt.Errorf("%s blinded block is nil", block.Version.String())
			}
			sb := &apiv1deneb.SignedBlindedBeaconBlock{
				Message:   block.DenebBlinded,
				Signature: sig,
			}
			r, _ = sb.HashTreeRoot()
		case spec.DataVersionElectra:
			if block.ElectraBlinded == nil {
				return fmt.Errorf("%s blinded block is nil", block.Version.String())
			}
			sb := &apiv1electra.SignedBlindedBeaconBlock{
				Message:   block.ElectraBlinded,
				Signature: sig,
			}
			r, _ = sb.HashTreeRoot()
		case spec.DataVersionFulu:
			if block.FuluBlinded == nil {
				return fmt.Errorf("%s blinded block is nil", block.Version.String())
			}
			sb := &apiv1electra.SignedBlindedBeaconBlock{
				Message:   block.FuluBlinded,
				Signature: sig,
			}
			r, _ = sb.HashTreeRoot()
		default:
			return fmt.Errorf("unknown blinded block version %d", block.Version)
		}
	} else {
		// Handle full blocks
		switch block.Version {
		case spec.DataVersionCapella:
			if block.Capella == nil {
				return fmt.Errorf("%s block is nil", block.Version.String())
			}
			sb := &capella.SignedBeaconBlock{
				Message:   block.Capella,
				Signature: sig,
			}
			r, _ = sb.HashTreeRoot()
		case spec.DataVersionDeneb:
			if block.Deneb == nil {
				return fmt.Errorf("%s block contents is nil", block.Version.String())
			}
			if block.Deneb.Block == nil {
				return fmt.Errorf("%s block is nil", block.Version.String())
			}
			sb := &apiv1deneb.SignedBlockContents{
				SignedBlock: &deneb.SignedBeaconBlock{
					Message:   block.Deneb.Block,
					Signature: sig,
				},
				KZGProofs: block.Deneb.KZGProofs,
				Blobs:     block.Deneb.Blobs,
			}
			r, _ = sb.HashTreeRoot()
		case spec.DataVersionElectra:
			if block.Electra == nil {
				return fmt.Errorf("%s block contents is nil", block.Version.String())
			}
			if block.Electra.Block == nil {
				return fmt.Errorf("%s block is nil", block.Version.String())
			}
			sb := &apiv1electra.SignedBlockContents{
				SignedBlock: &electra.SignedBeaconBlock{
					Message:   block.Electra.Block,
					Signature: sig,
				},
				KZGProofs: block.Electra.KZGProofs,
				Blobs:     block.Electra.Blobs,
			}
			r, _ = sb.HashTreeRoot()
		case spec.DataVersionFulu:
			if block.Fulu == nil {
				return fmt.Errorf("%s block contents is nil", block.Version.String())
			}
			if block.Fulu.Block == nil {
				return fmt.Errorf("%s block is nil", block.Version.String())
			}
			sb := &apiv1fulu.SignedBlockContents{
				SignedBlock: &electra.SignedBeaconBlock{
					Message:   block.Fulu.Block,
					Signature: sig,
				},
				KZGProofs: block.Fulu.KZGProofs,
				Blobs:     block.Fulu.Blobs,
			}
			r, _ = sb.HashTreeRoot()
		default:
			return types.WrapError(types.UnknownBlockVersionErrorCode, fmt.Errorf("unknown block version %d", block.Version))
		}
	}

	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// IsAggregator returns true if the validator is selected as an aggregator
func (bn *TestingBeaconNode) IsAggregator(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeeLength uint64, slotSig []byte) bool {
	// Simple mock: always return true for testing
	// In production, this would check the selection proof against the committee modulo
	return true
}

// GetAggregateAttestation returns the aggregate attestation for the given slot and committee
func (bn *TestingBeaconNode) GetAggregateAttestation(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (ssz.Marshaler, error) {
	version := VersionBySlot(slot)
	if version == spec.DataVersionElectra {
		return TestingElectraAggregateAndProof(TestingValidatorIndex).Aggregate, nil
	} else {
		return TestingPhase0AggregateAndProof(TestingValidatorIndex).Aggregate, nil
	}
}

// SubmitAggregateSelectionProof returns an AggregateAndProof object
// Deprecated: Use IsAggregator and GetAggregateAttestation instead. Kept for backward compatibility.
func (bn *TestingBeaconNode) SubmitAggregateSelectionProof(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeeLength uint64, index phase0.ValidatorIndex, slotSig []byte) (ssz.Marshaler, spec.DataVersion, error) {
	version := VersionBySlot(slot)
	return TestingAggregateAndProofV(version, TestingValidatorIndex), version, nil
}

// SubmitSignedAggregateSelectionProof broadcasts a signed aggregator msg
func (bn *TestingBeaconNode) SubmitSignedAggregateSelectionProof(msg *spec.VersionedSignedAggregateAndProof) error {
	var root [32]byte

	switch msg.Version {
	case spec.DataVersionPhase0:
		root, _ = msg.Phase0.HashTreeRoot()
	case spec.DataVersionAltair:
		root, _ = msg.Altair.HashTreeRoot()
	case spec.DataVersionBellatrix:
		root, _ = msg.Bellatrix.HashTreeRoot()
	case spec.DataVersionCapella:
		root, _ = msg.Capella.HashTreeRoot()
	case spec.DataVersionDeneb:
		root, _ = msg.Deneb.HashTreeRoot()
	case spec.DataVersionElectra:
		root, _ = msg.Electra.HashTreeRoot()
	case spec.DataVersionFulu:
		root, _ = msg.Fulu.HashTreeRoot()
	default:
		panic("unsupported version")
	}

	bn.BroadcastedRoots = append(bn.BroadcastedRoots, root)
	return nil
}

// GetSyncMessageBlockRoot returns beacon block root for sync committee
func (bn *TestingBeaconNode) GetSyncMessageBlockRoot(slot phase0.Slot) (phase0.Root, spec.DataVersion, error) {
	return TestingSyncCommitteeBlockRoot, spec.DataVersionPhase0, nil
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
	if len(bn.syncCommitteeAggregatorRoots) != 0 {
		if val, found := bn.syncCommitteeAggregatorRoots[hex.EncodeToString(proof)]; found {
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
	// return index / subnetSize

	// For testing, we use a simplified approach where subnet ID equals the index
	// This means validator at index 0 is in subnet 0, index 1 in subnet 1, etc.
	// In production, multiple validators share the same subnet (128 validators per subnet)
	return uint64(index)
}

// GetSyncCommitteeContribution returns
func (bn *TestingBeaconNode) GetSyncCommitteeContribution(slot phase0.Slot, selectionProofs []phase0.BLSSignature, subnetIDs []uint64) (ssz.Marshaler, spec.DataVersion, error) {
	return &TestingContributionsData, spec.DataVersionBellatrix, nil
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

func (bn *TestingBeaconNode) DataVersion(epoch phase0.Epoch) spec.DataVersion {
	return VersionByEpoch(epoch)
}
