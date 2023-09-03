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
	// Duty max size is
	// 			8 + 48 + 6*8 + 13*8 = 208 ~= 2^8
	Duty    Duty
	Version spec.DataVersion
	// PreConsensusJustifications max size is
	//			13*SignedPartialSignatureMessage(2^16) ~= 2^20
	PreConsensusJustifications []*SignedPartialSignatureMessage `ssz-max:"13"`
	// DataSSZ has max size as following
	// Biggest object is a full beacon block
	// BeaconBlock is 2*32+2*8 + BeaconBlockBody
	// BeaconBlockBody is
	//			96 + ETH1Data(2*32+8) + 32 +
	//			16*ProposerSlashing(2*SignedBeaconBlockHeader(96 + 3*32 + 2*8)) +
	//			2*AttesterSlashing(2*IndexedAttestation(2048*8 + 96 + AttestationData(2*8 + 32 + 2*(8+32)))) +
	//			128*Attestation(2048*8 + 96 + AttestationData(2*8 + 32 + 2*(8+32))) +
	//			16*Deposit(33*32 + 48 + 32 + 8 + 96) +
	//			16*SignedVoluntaryExit(96 + 2*8) +
	//			SyncAggregate(64 + 96) +
	//			ExecutionPayload(32 + 20 + 2*32 + 256 + 32 + 4*8 + 3*32 + 1048576*1073741824)
	// = 2^21(everything but transactions) + 2^50 (transaction list)
	// We do not need to support such a big DataSSZ size as 2^50 represents 1000X the actual block gas limit
	// Current 30M gas limit produces 30M / 16 (call data cost) = 1,875,000 bytes (https://eips.ethereum.org/EIPS/eip-4488)
	// we can upper limit transactions to 2^21, together with the rest of the data 2*2^21 = 2^22 = 4,194,304 bytes
	// Exmplanation on why transaction sizes are so big https://github.com/ethereum/consensus-specs/pull/2686
	DataSSZ []byte `ssz-max:"4194304"` // 2^22
}

func (cid *ConsensusData) Validate() error {
	switch cid.Duty.Type {
	case BNRoleAttester:
		if _, err := cid.GetAttestationData(); err != nil {
			return err
		}
		if len(cid.PreConsensusJustifications) > 0 {
			return errors.New("attester invalid justifications")
		}
	case BNRoleAggregator:
		if _, err := cid.GetAggregateAndProof(); err != nil {
			return err
		}
	case BNRoleProposer:
		var err1, err2 error
		_, _, err1 = cid.GetBlockData()
		_, _, err2 = cid.GetBlindedBlockData()

		if err1 != nil && err2 != nil {
			return err1
		}
		if err1 == nil && err2 == nil {
			return errors.New("no beacon data")
		}
	case BNRoleSyncCommittee:
		if len(cid.PreConsensusJustifications) > 0 {
			return errors.New("sync committee invalid justifications")
		}
		if _, err := cid.GetSyncCommitteeBlockRoot(); err != nil {
			return err
		}
		return nil
	case BNRoleSyncCommitteeContribution:
		if _, err := cid.GetSyncCommitteeContributions(); err != nil {
			return err
		}
	case BNRoleValidatorRegistration:
		return errors.New("validator registration has no consensus data")
	case BNRoleVoluntaryExit:
		return errors.New("voluntary exit has no consensus data")
	default:
		return errors.New("unknown duty role")
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

// GetBlockData ISSUE 221: GetBlockData/GetBlindedBlockData return versioned block only
func (ci *ConsensusData) GetBlockData() (*spec.VersionedBeaconBlock, ssz.HashRoot, error) {
	switch ci.Version {
	case spec.DataVersionBellatrix:
		ret := &bellatrix.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &spec.VersionedBeaconBlock{Bellatrix: ret, Version: ci.Version}, ret, nil
	case spec.DataVersionCapella:
		ret := &capella.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &spec.VersionedBeaconBlock{Capella: ret, Version: ci.Version}, ret, nil
	default:
		return nil, nil, errors.Errorf("unknown block version %s", ci.Version.String())
	}
}

// GetBlindedBlockData ISSUE 221: GetBlockData/GetBlindedBlockData return versioned block only
func (ci *ConsensusData) GetBlindedBlockData() (*api.VersionedBlindedBeaconBlock, ssz.HashRoot, error) {
	switch ci.Version {
	case spec.DataVersionBellatrix:
		ret := &apiv1bellatrix.BlindedBeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &api.VersionedBlindedBeaconBlock{Bellatrix: ret, Version: ci.Version}, ret, nil
	case spec.DataVersionCapella:
		ret := &apiv1capella.BlindedBeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &api.VersionedBlindedBeaconBlock{Capella: ret, Version: ci.Version}, ret, nil
	default:
		return nil, nil, errors.Errorf("unknown blinded block version %s", ci.Version.String())
	}
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
