package types

import (
	"encoding/hex"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
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

			for k, v := range cd.SyncCommitteeContribution {
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
	return json.Marshal(cd)
}

func (cd *ConsensusData) Decode(data []byte) error {
	return json.Unmarshal(data, &cd)
}

type ConsensusInput struct {
	Duty *Duty
	// TODO: determine real ssz-max. the current ssz-max calculated for the altair.BeaconBlock and not bellatrix.
	// bellatrix includes the transactions and th extra data inside the ExecutionPayload
	// Transactions  []Transaction `ssz-max:"1073741824,1048576"`
	Data []byte `ssz-max:"387068"`
}
