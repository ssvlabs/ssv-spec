package types

import (
	"github.com/attestantio/go-eth2-client/api"
	bellatrix2 "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	capella2 "github.com/attestantio/go-eth2-client/api/v1/capella"
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
		switch cid.Version {
		case spec.DataVersionBellatrix:
			_, err1 = cid.GetBellatrixBlockData()
			_, err2 = cid.GetBellatrixBlindedBlockData()
		case spec.DataVersionCapella:
			_, err1 = cid.GetCapellaBlockData()
			_, err2 = cid.GetCapellaBlindedBlockData()
		default:
			return errors.New("invalid block data")
		}

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

func (ci *ConsensusData) GetBellatrixBlockData() (*bellatrix.BeaconBlock, error) {
	ret := &bellatrix.BeaconBlock{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ssz")
	}
	return ret, nil
}

// DecidedBlindedBlock returns true if decided value has a blinded block, false if regular block
// WARNING!! should be called after decided only
func (ci *ConsensusData) DecidedBlindedBlock() bool {
	switch ci.Version {
	case spec.DataVersionBellatrix:
		_, err := ci.GetBellatrixBlindedBlockData()
		return err == nil
	case spec.DataVersionCapella:
		_, err := ci.GetCapellaBlindedBlockData()
		return err == nil
	default:
		return false
	}
}

func (ci *ConsensusData) GetBlockRoot() (ssz.HashRoot, error) {
	switch ci.Version {
	case spec.DataVersionBellatrix:
		blk, err := ci.GetBellatrixBlindedBlockData()
		if err == nil { // if no error, is blinded block
			return blk, nil
		}
		return ci.GetBellatrixBlindedBlockData()
	case spec.DataVersionCapella:
		blk, err := ci.GetCapellaBlindedBlockData()
		if err == nil { // if no error, is blinded block
			return blk, nil
		}
		return ci.GetCapellaBlockData()
	default:
		return nil, errors.New("not supported version")
	}
}

func (ci *ConsensusData) GetVersionedBlock(sig phase0.BLSSignature) (*spec.VersionedSignedBeaconBlock, error) {
	switch ci.Version {
	case spec.DataVersionBellatrix:
		data, err := ci.GetBellatrixBlockData()
		if err != nil {
			return nil, errors.Wrap(err, "could not get bellatrix block")
		}

		blk := &bellatrix.SignedBeaconBlock{
			Message:   data,
			Signature: sig,
		}
		return &spec.VersionedSignedBeaconBlock{
			Version:   spec.DataVersionBellatrix,
			Bellatrix: blk,
		}, nil
	case spec.DataVersionCapella:
		data, err := ci.GetCapellaBlockData()
		if err != nil {
			return nil, errors.Wrap(err, "could not get vapella block")
		}

		blk := &capella.SignedBeaconBlock{
			Message:   data,
			Signature: sig,
		}
		return &spec.VersionedSignedBeaconBlock{
			Version: spec.DataVersionCapella,
			Capella: blk,
		}, nil
	default:
		return nil, errors.New("not supported version")
	}
}

func (ci *ConsensusData) GetVersionedBlindedBlock(sig phase0.BLSSignature) (*api.VersionedSignedBlindedBeaconBlock, error) {
	switch ci.Version {
	case spec.DataVersionBellatrix:
		data, err := ci.GetBellatrixBlindedBlockData()
		if err != nil {
			return nil, errors.Wrap(err, "could not get bellatrix block")
		}

		blk := &bellatrix2.SignedBlindedBeaconBlock{
			Message:   data,
			Signature: sig,
		}
		return &api.VersionedSignedBlindedBeaconBlock{
			Version:   spec.DataVersionBellatrix,
			Bellatrix: blk,
		}, nil
	case spec.DataVersionCapella:
		data, err := ci.GetCapellaBlindedBlockData()
		if err != nil {
			return nil, errors.Wrap(err, "could not get vapella block")
		}

		blk := &capella2.SignedBlindedBeaconBlock{
			Message:   data,
			Signature: sig,
		}
		return &api.VersionedSignedBlindedBeaconBlock{
			Version: spec.DataVersionCapella,
			Capella: blk,
		}, nil
	default:
		return nil, errors.New("not supported version")
	}
}

func (ci *ConsensusData) GetBellatrixBlindedBlockData() (*bellatrix2.BlindedBeaconBlock, error) {
	ret := &bellatrix2.BlindedBeaconBlock{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ssz")
	}
	return ret, nil
}

func (ci *ConsensusData) GetCapellaBlockData() (*capella.BeaconBlock, error) {
	ret := &capella.BeaconBlock{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ssz")
	}
	return ret, nil
}

func (ci *ConsensusData) GetCapellaBlindedBlockData() (*capella2.BlindedBeaconBlock, error) {
	ret := &capella2.BlindedBeaconBlock{}
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
