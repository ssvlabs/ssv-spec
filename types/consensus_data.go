package types

import (
	"github.com/attestantio/go-eth2-client/api"
	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
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
		return nil, ssz.ErrListTooBigFn("ConsensusData.SyncCommitteeContribution", size, 13)
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

// ConsensusData holds all relevant duty and data Decided on by consensus
type ConsensusData struct {
	Duty                       Duty
	Version                    spec.DataVersion
	PreConsensusJustifications []*SignedPartialSignatureMessage `ssz-max:"13"`
	DataSSZ                    []byte                           `ssz-max:"1073807360"` // 2^30+2^16 (considering max block size 2^30)
}

func (cid *ConsensusData) Validate() error {
	switch cid.Duty.Type {
	case BNRoleAttester:
		if _, err := cid.GetAttestationData(); err != nil {
			return err
		}
	case BNRoleAggregator:
		if _, err := cid.GetAggregateAndProof(); err != nil {
			return err
		}
	case BNRoleProposer:
		var err1, err2 error
		_, err1 = cid.GetBlockData()
		_, err2 = cid.GetBlindedBlockData()

		if err1 != nil && err2 != nil {
			return err1
		}
		if err1 == nil && err2 == nil {
			return errors.New("no beacon data")
		}
	case BNRoleSyncCommittee:
		return nil
	case BNRoleSyncCommitteeContribution:
		if _, err := cid.GetSyncCommitteeContributions(); err != nil {
			return err
		}
	}
	return nil
}

func (ci *ConsensusData) GetAttestationData() (*phase0.AttestationData, error) {
	ret := &phase0.AttestationData{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ssz")
	}
	return ret, nil
}

func (ci *ConsensusData) GetBlockData() (*spec.VersionedBeaconBlock, error) {
	switch ci.Version {
	case spec.DataVersionPhase0:
		ret := &phase0.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &spec.VersionedBeaconBlock{Phase0: ret, Version: spec.DataVersionPhase0}, nil
	case spec.DataVersionAltair:
		ret := &altair.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &spec.VersionedBeaconBlock{Altair: ret, Version: spec.DataVersionAltair}, nil
	case spec.DataVersionBellatrix:
		ret := &bellatrix.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &spec.VersionedBeaconBlock{Bellatrix: ret, Version: spec.DataVersionBellatrix}, nil
	case spec.DataVersionCapella:
		ret := &capella.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &spec.VersionedBeaconBlock{Capella: ret, Version: spec.DataVersionCapella}, nil
	default:
		return nil, errors.New("unknown beacon block version")
	}
}

func (ci *ConsensusData) GetBlockHashRoot() (ssz.HashRoot, error) {
	switch ci.Version {
	case spec.DataVersionPhase0:
		ret := &phase0.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return ret, nil
	case spec.DataVersionAltair:
		ret := &altair.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return ret, nil
	case spec.DataVersionBellatrix:
		ret := &bellatrix.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return ret, nil
	case spec.DataVersionCapella:
		ret := &capella.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return ret, nil
	default:
		return nil, errors.New("unknown beacon block version")
	}
}

func (ci *ConsensusData) GetBlindedBlockData() (*api.VersionedBlindedBeaconBlock, error) {
	switch ci.Version {
	case spec.DataVersionBellatrix:
		ret := &apiv1bellatrix.BlindedBeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &api.VersionedBlindedBeaconBlock{Bellatrix: ret, Version: spec.DataVersionBellatrix}, nil
	case spec.DataVersionCapella:
		ret := &apiv1capella.BlindedBeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &api.VersionedBlindedBeaconBlock{Capella: ret, Version: spec.DataVersionCapella}, nil
	default:
		return nil, errors.New("unknown blinded beacon block version")
	}
}

func (ci *ConsensusData) GetBlindedBlockHashRoot() (ssz.HashRoot, error) {
	switch ci.Version {
	case spec.DataVersionBellatrix:
		ret := &apiv1bellatrix.BlindedBeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return ret, nil
	case spec.DataVersionCapella:
		ret := &apiv1capella.BlindedBeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return ret, nil
	}

	return nil, errors.New("unknown blinded beacon block version")
}

func (ci *ConsensusData) GetBellatrixBlindedBlockData() (*apiv1bellatrix.BlindedBeaconBlock, error) {
	ret := &apiv1bellatrix.BlindedBeaconBlock{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ssz")
	}
	return ret, nil
}

func (ci *ConsensusData) GetAggregateAndProof() (*phase0.AggregateAndProof, error) {
	ret := &phase0.AggregateAndProof{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ssz")
	}
	return ret, nil
}

func (ci *ConsensusData) GetSyncCommitteeBlockRoot() (phase0.Root, error) {
	ret := SSZ32Bytes{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return phase0.Root{}, errors.Wrap(err, "could not unmarshal ssz")
	}
	return phase0.Root(ret), nil
}

func (ci *ConsensusData) GetSyncCommitteeContributions() (Contributions, error) {
	ret := Contributions{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ssz")
	}
	return ret, nil
}

func (cid *ConsensusData) Encode() ([]byte, error) {
	return cid.MarshalSSZ()
}

func (cid *ConsensusData) Decode(data []byte) error {
	return cid.UnmarshalSSZ(data)
}
