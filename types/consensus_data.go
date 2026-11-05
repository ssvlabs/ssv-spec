package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types/gloas"
)

type Contribution struct {
	SelectionProofSig [96]byte `ssz-size:"96"`
	Contribution      altair.SyncCommitteeContribution
}

// Contributions --
type Contributions []*Contribution

func (c *Contributions) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(c)
}

func (c *Contributions) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(c)
}

func (c *Contributions) HashTreeRootWith(hh ssz.HashWalker) error {
	// taken from https://github.com/prysmaticlabs/prysm/blob/develop/encoding/ssz/htrutils.go#L97-L119
	subIndx := hh.Index()
	num := uint64(len(*c))
	if num > 13 {
		return ssz.ErrIncorrectListSize
	}
	for _, elem := range *c {
		{
			if err := elem.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
	}
	hh.MerkleizeWithMixin(subIndx, num, 13)
	return nil
}

// UnmarshalSSZ --
func (c *Contributions) UnmarshalSSZ(buf []byte) error {
	num, err := ssz.DecodeDynamicLength(buf, 13)
	if err != nil {
		return err
	}
	*c = make(Contributions, num)

	return ssz.UnmarshalDynamic(buf, num, func(indx int, buf []byte) (err error) {
		if (*c)[indx] == nil {
			(*c)[indx] = new(Contribution)
		}
		if err = (*c)[indx].UnmarshalSSZ(buf); err != nil {
			return err
		}
		return nil
	})
}

// MarshalSSZTo --
func (c *Contributions) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	if size := len(*c); size > 13 {
		return nil, ssz.ErrListTooBigFn("ValidatorConsensusData.SyncCommitteeContribution", size, 13)
	}

	offset := 4 * len(*c)
	for ii := 0; ii < len(*c); ii++ {
		dst = ssz.WriteOffset(dst, offset)
		offset += (*c)[ii].SizeSSZ()
	}

	for ii := 0; ii < len(*c); ii++ {
		if dst, err = (*c)[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}
	return dst, nil
}

// MarshalSSZ --
func (c *Contributions) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

// SizeSSZ returns the size of the serialized object.
func (c Contributions) SizeSSZ() int {
	size := 0
	for _, elem := range c {
		size += 4
		size += elem.SizeSSZ()
	}
	return size
}

// BeaconVote is the value the CommitteeRunner agrees on (SIP #94 §2): the head block root, the FFG checkpoints,
// and the BN-supplied AttestationData.Index (120-byte SSZ). The index is the fork-choice payload status (0 =
// EMPTY, 1 = FULL; 0 for a same-slot attestation) and part of the signed attestation root, so it travels
// through consensus.
type BeaconVote struct {
	BlockRoot            phase0.Root `ssz-size:"32"`
	Source               *phase0.Checkpoint
	Target               *phase0.Checkpoint
	AttestationDataIndex phase0.CommitteeIndex // copied from AttestationData.Index (0 or 1)
}

// Encode the BeaconVote object
func (b *BeaconVote) Encode() ([]byte, error) {
	return b.MarshalSSZ()
}

// Decode the BeaconVote object
func (b *BeaconVote) Decode(data []byte) error {
	return b.UnmarshalSSZ(data)
}

// Validate checks the following rules:
//   - Source and Target checkpoints must be non-nil
//   - AttestationDataIndex must be 0 or 1 (the payload status; SIP #94 §2)
//   - Source.Epoch must be strictly less than Target.Epoch
func (b *BeaconVote) Validate() error {
	if b == nil {
		return NewError(BeaconVoteNilCheckpointErrorCode, "nil beacon vote")
	}
	if b.Source == nil || b.Target == nil {
		return NewError(BeaconVoteNilCheckpointErrorCode, "nil source or target checkpoint")
	}
	if b.AttestationDataIndex > 1 {
		return NewError(BeaconVoteInvalidIndexErrorCode,
			fmt.Sprintf("attestation data index %d must be 0 or 1", b.AttestationDataIndex))
	}
	if b.Source.Epoch >= b.Target.Epoch {
		return NewError(AttestationSourceNotLessThanTargetErrorCode, "attestation data source >= target")
	}
	return nil
}

// ProposerConsensusData holds all relevant data about proposer duty for consensus
type ProposerConsensusData struct {
	// Duty max size is
	// 			8 + 48 + 6*8 + 13*8 + 1 = 209
	Duty ValidatorDuty
	// Version is the proposal's fork; it must be Gloas (see GetBlockData).
	Version spec.DataVersion
	// DataSSZ is the SSZ-encoded gloas.GloasProposalData: the block plus the §6 payload_root (SIP #94 §4). A Gloas
	// block carries a payload bid rather than the payload, so it stays far below this bound (2^23, kept from the
	// pre-Gloas BlockContents sizing).
	DataSSZ []byte `ssz-max:"8388608"`
}

// versionJSON is the JSON codec for a consensus-data Version: "gloas" for the SIP #94 placeholder, the upstream
// fork string for other known versions, and the bare number for any other out-of-enum value (on which
// spec.DataVersion.MarshalJSON panics). Decoding accepts the string (any case), the number, and a missing/null
// value (as version 0).
type versionJSON spec.DataVersion

func (v versionJSON) MarshalJSON() ([]byte, error) {
	dv := spec.DataVersion(v)
	switch {
	case dv == gloas.DataVersionGloas:
		return json.Marshal("gloas")
	case dv.String() == "unknown": // out-of-enum: spec.DataVersion.MarshalJSON would panic
		return json.Marshal(uint64(dv))
	default:
		return dv.MarshalJSON()
	}
}

func (v *versionJSON) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*v = 0
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if strings.EqualFold(s, "gloas") {
			*v = versionJSON(gloas.DataVersionGloas)
			return nil
		}
		var dv spec.DataVersion
		if err := dv.UnmarshalJSON(data); err != nil {
			return err
		}
		*v = versionJSON(dv)
		return nil
	}
	var n uint64
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*v = versionJSON(n)
	return nil
}

// MarshalJSON/UnmarshalJSON write Version through versionJSON, listing the fields explicitly to keep the key
// order (Duty, Version, DataSSZ); a new struct field must be added to both (TestProposerConsensusDataJSONFieldsInSync).
func (cd *ProposerConsensusData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Duty    ValidatorDuty
		Version versionJSON
		DataSSZ []byte
	}{cd.Duty, versionJSON(cd.Version), cd.DataSSZ})
}

func (cd *ProposerConsensusData) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Duty    ValidatorDuty
		Version versionJSON
		DataSSZ []byte
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	cd.Duty, cd.Version, cd.DataSSZ = aux.Duty, spec.DataVersion(aux.Version), aux.DataSSZ
	return nil
}

func (cd *ProposerConsensusData) Validate() error {
	if cd.Duty.Type != BNRoleProposer {
		return NewError(UnknownDutyRoleDataErrorCode, "unknown duty role")
	}
	_, err := cd.GetBlockData()
	return err
}

// GetBlockData decodes the proposal: the Gloas block and the §6 payload_root carried with it (SIP #94 §4).
func (cd *ProposerConsensusData) GetBlockData() (*gloas.GloasProposalData, error) {
	if cd.Version != gloas.DataVersionGloas {
		return nil, WrapError(UnknownBlockVersionErrorCode, fmt.Errorf("unknown block version %d", cd.Version))
	}
	proposalData, err := gloas.DecodeGloasProposalData(cd.DataSSZ)
	if err != nil {
		return nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal gloas proposal data: %w", err))
	}
	return proposalData, nil
}

func (cd *ProposerConsensusData) Encode() ([]byte, error) {
	return cd.MarshalSSZ()
}

func (cd *ProposerConsensusData) Decode(data []byte) error {
	return cd.UnmarshalSSZ(data)
}

// AssignedAggregator represents a validator that has been assigned as an aggregator or sync committee contributor
type AssignedAggregator struct {
	ValidatorIndex phase0.ValidatorIndex
	SelectionProof phase0.BLSSignature `ssz-size:"96"`
	CommitteeIndex uint64
}

// AggregatorCommitteeConsensusData is the consensus data for the aggregator committee runner
type AggregatorCommitteeConsensusData struct {
	Version spec.DataVersion

	// Aggregator duties
	Aggregators []AssignedAggregator `ssz-max:"3000"` // For a maximum of 3k validators per committee
	// AggregatorsCommitteeIndexes is a list of committee indexes used by the above aggregators
	AggregatorsCommitteeIndexes []uint64 `ssz-max:"64"`
	// AggregatedAttestations is a list of aggregated attestations (SSZ bytes), one for each committee above
	AggregatedAttestations [][]byte `ssz-max:"64,131308"`

	// Sync Committee duties
	Contributors []AssignedAggregator `ssz-max:"2048"` // 512 * 4
	// SyncCommitteeContributions is a list of contributions, one for each subcommittee
	SyncCommitteeContributions []altair.SyncCommitteeContribution `ssz-max:"4"`
}

// MarshalJSON writes Version through versionJSON so a Gloas-stamped value is JSON-safe. Version is the
// struct's first field, so the embedded-alias order already matches.
func (a *AggregatorCommitteeConsensusData) MarshalJSON() ([]byte, error) {
	type alias AggregatorCommitteeConsensusData
	return json.Marshal(&struct {
		Version versionJSON
		*alias
	}{versionJSON(a.Version), (*alias)(a)})
}

func (a *AggregatorCommitteeConsensusData) UnmarshalJSON(data []byte) error {
	type alias AggregatorCommitteeConsensusData
	aux := &struct {
		Version versionJSON
		*alias
	}{alias: (*alias)(a)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	a.Version = spec.DataVersion(aux.Version)
	return nil
}

// Validate ensures the consensus data is internally consistent
func (a *AggregatorCommitteeConsensusData) Validate() error {
	if a.Version != gloas.DataVersionGloas {
		return WrapError(UnknownVersionErrorCode, fmt.Errorf("unknown version %d", a.Version))
	}

	// Ensure at least one validator
	if len(a.Aggregators) == 0 && len(a.Contributors) == 0 {
		return NewError(AggCommConsensusDataNoValidatorErrorCode, "no validators assigned to aggregator committee or sync committee")
	}

	// Aggregators validation

	// Ensure there is exactly one aggregated attestation per committee index
	if len(a.AggregatorsCommitteeIndexes) != len(a.AggregatedAttestations) {
		return NewError(AggCommAggCommIdxCntMismatchErrorCode, "committee indexes and attestations count mismatch")
	}

	// Validate equal set (AggregatorsCommitteeIndexes vs. Aggregators.CommitteeIndex)
	allowedAggCommittees := make(map[uint64]struct{}, len(a.AggregatorsCommitteeIndexes))
	for _, idx := range a.AggregatorsCommitteeIndexes {
		// Duplicates are not allowed
		if _, dup := allowedAggCommittees[idx]; dup {
			return NewError(AggCommDuplicatedCommIdxErrorCode, "duplicate index in AggregatorsCommitteeIndexes")
		}
		allowedAggCommittees[idx] = struct{}{}
	}
	usedAggCommittees := make(map[uint64]struct{}, len(a.AggregatorsCommitteeIndexes))
	for _, agg := range a.Aggregators {
		// Check it exists in allowed
		if _, ok := allowedAggCommittees[agg.CommitteeIndex]; !ok {
			return NewError(AggCommCommIdxMismatchErrorCode, "aggregator committee index not listed in AggregatorsCommitteeIndexes")
		}
		// Mark as used
		usedAggCommittees[agg.CommitteeIndex] = struct{}{}
	}
	// Ensure no committee index was left unused (no more than necessary)
	if len(usedAggCommittees) != len(allowedAggCommittees) {
		return NewError(AggCommUnusedCommIdxErrorCode, "leftover aggregator committee index not usedAggCommittees by any aggregator")
	}

	// Ensure attestation objects are decoded correctly. Gloas reuses the Electra attestation shape (SIP #94 §2).
	for _, attBytes := range a.AggregatedAttestations {
		att := &electra.Attestation{}
		if err := att.UnmarshalSSZ(attBytes); err != nil {
			return NewError(AggCommAttestationDecodingErrorCode, "failed to unmarshal attestation")
		}
	}

	// Sync committee contributors validation

	// Validate equal set (Contributors.CommitteeIndex vs. SyncCommitteeContributions.SubcommitteeIndex)
	allowedSCSubnets := make(map[uint64]struct{}, len(a.SyncCommitteeContributions))
	for _, contrib := range a.SyncCommitteeContributions {
		// Duplicates are not allowed
		if _, dup := allowedSCSubnets[contrib.SubcommitteeIndex]; dup {
			return NewError(AggCommSCCSubnetDuplicateErrorCode, "duplicate subcommittee index in SyncCommitteeContributions")
		}
		allowedSCSubnets[contrib.SubcommitteeIndex] = struct{}{}
	}
	usedSCSubnets := make(map[uint64]struct{}, len(a.SyncCommitteeContributions))
	for _, contributor := range a.Contributors {
		// Check it exists in allowed
		if _, ok := allowedSCSubnets[contributor.CommitteeIndex]; !ok {
			return NewError(AggCommSubnetNotInSCSubnetsErrorCode, "sync committee contributor subnet not listed in SyncCommitteeContributions")
		}
		// Mark as used
		usedSCSubnets[contributor.CommitteeIndex] = struct{}{}
	}
	// Ensure no subcommittee index was left unused (no more than necessary)
	if len(usedSCSubnets) != len(allowedSCSubnets) {
		return NewError(AggCommUnusedSubnetErrorCode, "leftover sync committee contributor subnet not used in SyncCommitteeContributions")
	}

	return nil
}

// Encode encodes the consensus data to SSZ
func (a *AggregatorCommitteeConsensusData) Encode() ([]byte, error) {
	return a.MarshalSSZ()
}

// Decode decodes the consensus data from SSZ
func (a *AggregatorCommitteeConsensusData) Decode(data []byte) error {
	return a.UnmarshalSSZ(data)
}

func GetAggregateAndProofHashRoot(aggProof *spec.VersionedAggregateAndProof) (ssz.HashRoot, error) {
	if aggProof.Version != gloas.DataVersionGloas {
		return nil, WrapError(UnknownVersionErrorCode, fmt.Errorf("unknown version %d", aggProof.Version))
	}
	// Gloas reuses the Electra aggregate-and-proof shape (SIP #94 §2); no Gloas field on the versioned wrapper.
	return aggProof.Electra, nil
}

// GetAggregateAndProofs returns all aggregate and proofs for the aggregator duties along with their hash roots
func (a *AggregatorCommitteeConsensusData) GetAggregateAndProofs() ([]*spec.VersionedAggregateAndProof, error) {
	if a.Version != gloas.DataVersionGloas {
		return nil, WrapError(UnknownVersionErrorCode, fmt.Errorf("unknown version %d", a.Version))
	}

	proofs := make([]*spec.VersionedAggregateAndProof, 0, len(a.Aggregators))

	for _, aggregator := range a.Aggregators {
		// Get index for validator in a.AggregatedAttestations
		foundIndex := -1
		for idx, committeeIndex := range a.AggregatorsCommitteeIndexes {
			if committeeIndex == aggregator.CommitteeIndex {
				foundIndex = idx
				break
			}
		}
		if foundIndex == -1 || foundIndex >= len(a.AggregatedAttestations) {
			return nil, NewError(AggCommCommIdxMismatchErrorCode, "aggregator committee index not found for attestation")
		}

		// Gloas reuses the Electra aggregate shape (SIP #94 §2).
		att := &electra.Attestation{}
		if err := att.UnmarshalSSZ(a.AggregatedAttestations[foundIndex]); err != nil {
			return nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("failed to unmarshal attestation: %w", err))
		}
		proofs = append(proofs, &spec.VersionedAggregateAndProof{
			Version: a.Version,
			Electra: &electra.AggregateAndProof{
				AggregatorIndex: aggregator.ValidatorIndex,
				Aggregate:       att,
				SelectionProof:  aggregator.SelectionProof,
			},
		})
	}

	return proofs, nil
}

// GetSyncCommitteeContributions returns the sync committee contributions
func (a *AggregatorCommitteeConsensusData) GetSyncCommitteeContributions() (Contributions, error) {

	contributions := make(Contributions, 0, len(a.Contributors))

	for _, contributor := range a.Contributors {

		// Find associated object in a.SyncCommitteeContributions
		foundIndex := -1
		for idx, contrib := range a.SyncCommitteeContributions {
			if contrib.SubcommitteeIndex == contributor.CommitteeIndex {
				foundIndex = idx
				break
			}
		}
		if foundIndex == -1 {
			return nil, NewError(AggCommSubnetNotInSCSubnetsErrorCode, "sync committee contributor subnet not found in SyncCommitteeContributions")
		}

		var sigBytes [96]byte
		copy(sigBytes[:], contributor.SelectionProof[:])
		contributions = append(contributions, &Contribution{
			SelectionProofSig: sigBytes,
			Contribution:      a.SyncCommitteeContributions[foundIndex],
		})
	}

	return contributions, nil
}
