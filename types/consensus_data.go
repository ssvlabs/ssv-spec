package types

import (
	"encoding/hex"
	"encoding/json"
	bellatrix2 "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type ContributionsMap map[phase0.BLSSignature]*altair.SyncCommitteeContribution

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
	BlindedBlockData       *bellatrix2.BlindedBeaconBlock
	AggregateAndProof      *phase0.AggregateAndProof
	SyncCommitteeBlockRoot phase0.Root
	// SyncCommitteeContribution map holds as key the selection proof for the contribution
	SyncCommitteeContribution ContributionsMap
}

func (cid *ConsensusData) ValidateForDuty(role BeaconRole) error {
	if cid.Duty == nil {
		return errors.New("duty is nil")
	}

	if role == BNRoleAttester && cid.AttestationData == nil {
		return errors.New("attestation data is nil")
	}

	if role == BNRoleAggregator && cid.AggregateAndProof == nil {
		return errors.New("aggregate and proof data is nil")
	}

	if role == BNRoleProposer {
		if cid.BlockData == nil && cid.BlindedBlockData == nil {
			return errors.New("block data is nil")
		}

		if cid.BlockData != nil && cid.BlindedBlockData != nil {
			return errors.New("block and blinded block data are both != nil")
		}
	}

	if role == BNRoleSyncCommitteeContribution && len(cid.SyncCommitteeContribution) == 0 {
		return errors.New("sync committee contribution data is nil")
	}

	return nil
}

func (cid *ConsensusData) Encode() ([]byte, error) {
	return json.Marshal(cid)
}

func (cid *ConsensusData) Decode(data []byte) error {
	return json.Unmarshal(data, &cid)
}
