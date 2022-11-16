package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
	"sort"
)

type ContributionsMap map[phase0.BLSSignature]*altair.SyncCommitteeContribution

type contributionEntries struct {
	SyncCommitteeContribution []*ContributionEntry `ssz-max:"13"`
}

type ContributionEntry struct {
	Sig   phase0.BLSSignature `ssz-size:"96"`
	Contr *altair.SyncCommitteeContribution
}

func (cm *ContributionsMap) MarshalJSON() ([]byte, error) {
	m := make(map[string]*altair.SyncCommitteeContribution)
	for k, v := range *cm {
		m[hex.EncodeToString(k[:])] = v
	}
	return json.Marshal(m)
}

func (cm *ContributionsMap) UnmarshalJSON(input []byte) error {
	m := make(map[string]*altair.SyncCommitteeContribution)
	if err := json.Unmarshal(input, &m); err != nil {
		return err
	}

	if *cm == nil {
		*cm = ContributionsMap{}
	}

	for k, v := range m {
		byts, err := hex.DecodeString(k)
		if err != nil {
			return err
		}

		blSig := phase0.BLSSignature{}
		copy(blSig[:], byts)
		(*cm)[blSig] = v
	}
	return nil
}

// ConsensusData holds all relevant duty and data Decided on by consensus
type ConsensusData struct {
	Duty                   *Duty
	AttestationData        *phase0.AttestationData
	BlockData              *bellatrix.BeaconBlock
	AggregateAndProof      *phase0.AggregateAndProof
	SyncCommitteeBlockRoot phase0.Root
	// SyncCommitteeContribution map holds as key the selection proof for the contribution
	SyncCommitteeContribution ContributionsMap
}

func (cd *ConsensusData) toConsensusInput() (*ConsensusInput, error) {
	var marshalSSZ []byte
	var err error
	if cd.Duty == nil {
		return nil, errors.New("could not marshal consensus data, duty is nil")
	}
	switch cd.Duty.Type {
	case BNRoleAttester:
		if cd.AttestationData == nil {
			return nil, errors.New("could not marshal consensus data, attestation data is nil")
		}
		marshalSSZ, err = cd.AttestationData.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	case BNRoleAggregator:
		if cd.AggregateAndProof == nil {
			return nil, errors.New("could not marshal consensus data, aggregate and proof is nil")
		}
		marshalSSZ, err = cd.AggregateAndProof.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	case BNRoleProposer:
		if cd.BlockData == nil {
			return nil, errors.New("could not marshal consensus data, block data is nil")
		}
		marshalSSZ, err = cd.BlockData.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	case BNRoleSyncCommittee:
		marshalSSZ = append(marshalSSZ, cd.SyncCommitteeBlockRoot[:]...)
	case BNRoleSyncCommitteeContribution:
		var ce contributionEntries
		if len(cd.SyncCommitteeContribution) > 0 {
			ce.SyncCommitteeContribution = make([]*ContributionEntry, 0, len(cd.SyncCommitteeContribution))

			keys := make([]phase0.BLSSignature, 0, len(cd.SyncCommitteeContribution))
			for k := range cd.SyncCommitteeContribution {
				keys = append(keys, k)
			}
			sort.Slice(keys, func(i, j int) bool {
				switch bytes.Compare(keys[i][:], keys[j][:]) {
				case -1:
					return true
				case 0, 1:
					return false
				default:
					return false
				}
			})

			for _, k := range keys {
				v := (cd.SyncCommitteeContribution)[k]
				ce.SyncCommitteeContribution = append(ce.SyncCommitteeContribution, &ContributionEntry{
					Sig:   k,
					Contr: v,
				})
			}
		}
		marshalSSZ, err = ce.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown role")
	}

	return &ConsensusInput{
		Duty: cd.Duty,
		Data: marshalSSZ,
	}, nil
}

func (cd *ConsensusData) Encode() ([]byte, error) {
	return cd.MarshalSSZ()
}

func (cd *ConsensusData) Decode(data []byte) error {
	return cd.UnmarshalSSZ(data)
}

func (cd *ConsensusData) MarshalSSZ() ([]byte, error) {
	ci, err := cd.toConsensusInput()
	if err != nil {
		return nil, err
	}

	return ssz.MarshalSSZ(ci)
}

// MarshalSSZTo ssz marshals the ConsensusData object to a target array
func (cd *ConsensusData) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	return nil, nil
}

// UnmarshalSSZ ssz unmarshals the ConsensusData object
func (cd *ConsensusData) UnmarshalSSZ(buf []byte) error {
	var cdSSZ ConsensusInput
	if err := cdSSZ.UnmarshalSSZ(buf); err != nil {
		return err
	}

	cd.Duty = cdSSZ.Duty

	var err error
	switch cd.Duty.Type {
	case BNRoleAttester:
		var attestationData phase0.AttestationData
		if err = attestationData.UnmarshalSSZ(cdSSZ.Data); err != nil {
			return err
		}
		cd.AttestationData = &attestationData
	case BNRoleAggregator:
		var aggregateAndProof phase0.AggregateAndProof
		if err = aggregateAndProof.UnmarshalSSZ(cdSSZ.Data); err != nil {
			return err
		}
		cd.AggregateAndProof = &aggregateAndProof
	case BNRoleProposer:
		var blockData bellatrix.BeaconBlock
		if err = blockData.UnmarshalSSZ(cdSSZ.Data); err != nil {
			return err
		}
		cd.BlockData = &blockData
	case BNRoleSyncCommittee:
		copy(cd.SyncCommitteeBlockRoot[:], cdSSZ.Data)
	case BNRoleSyncCommitteeContribution:
		var ce contributionEntries
		if err = ce.UnmarshalSSZ(cdSSZ.Data); err != nil {
			return err
		}
		var contributionMap ContributionsMap
		if len(ce.SyncCommitteeContribution) > 0 {
			contributionMap = make(ContributionsMap)
			for _, s := range ce.SyncCommitteeContribution {
				contributionMap[s.Sig] = s.Contr
			}
		}
		cd.SyncCommitteeContribution = contributionMap
	default:
		return errors.New("unknown role")
	}

	return nil
}

// SizeSSZ returns the ssz encoded size in bytes for the ConsensusData object
func (cd *ConsensusData) SizeSSZ() (size int) {
	size = 101

	switch cd.Duty.Type {
	case BNRoleAttester:
		size += cd.AttestationData.SizeSSZ()
	case BNRoleAggregator:
		size += cd.AggregateAndProof.SizeSSZ()
	case BNRoleProposer:
		size += cd.BlockData.SizeSSZ()
	case BNRoleSyncCommittee:
		size += len(cd.SyncCommitteeBlockRoot)
	case BNRoleSyncCommitteeContribution:
		sccSlice := make([]*ContributionEntry, 0, len(cd.SyncCommitteeContribution))

		for key, value := range cd.SyncCommitteeContribution {
			sccSlice = append(sccSlice, &ContributionEntry{
				Sig:   key,
				Contr: value,
			})
		}

		for ii := 0; ii < len(sccSlice); ii++ {
			size += 4
			size += sccSlice[ii].SizeSSZ()
		}
	}

	return
}

// HashTreeRoot ssz hashes the ConsensusData object
func (cd *ConsensusData) HashTreeRoot() ([32]byte, error) {
	ci, err := cd.toConsensusInput()
	if err != nil {
		return [32]byte{}, err
	}

	return ssz.HashWithDefaultHasher(ci)
}

type ConsensusInput struct {
	Duty *Duty
	// BeaconBlock max size
	// bellatrix includes the transactions and th extra data inside the ExecutionPayload
	// Transactions  []Transaction `ssz-max:"1073741824,1048576"`
	Data []byte `ssz-max:"1125899907230230"`
}
