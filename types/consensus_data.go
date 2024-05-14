package types

import (
	"github.com/attestantio/go-eth2-client/api"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
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
	Duty    BeaconDuty
	Version spec.DataVersion
	// PreConsensusJustifications max size is
	//			13*SignedPartialSignatureMessage(2^16) ~= 2^20
	PreConsensusJustifications []*PartialSignatureMessages `ssz-max:"13"`
	// DataSSZ has max size as following
	// Biggest object is a Deneb.BlockContents with size = 1120872 (everything but transactions) + 2^50 (transaction list)
	// We do not need to support such a big DataSSZ size as 2^50 represents 1000X the actual block gas limit
	// Upcoming 40M gas limit produces 40M / 16 (call data cost) = 2,500,000 bytes (https://eips.ethereum.org/EIPS/eip-4488)
	// Adding to the rest of the data, we have: 1,120,872 + 2,500,000  = 3,620,872 ~<= 2^21
	// Explanation on why transaction sizes are so big https://github.com/ethereum/consensus-specs/pull/2686
	// Python script for Deneb.BlockContents:
	// 		MaxTransactions = 2_500_000
	// 		Graffiti = 32
	// 		ETH1Data = 32 + 8 + 32
	// 		# ProposerSlashing
	// 		BeaconBlockHeader = 8 + 8 + 3 * 32
	// 		SignedBeaconBlockHeader = BeaconBlockHeader + 96
	// 		ProposerSlashing = 2 * SignedBeaconBlockHeader
	// 		# AttesterSlashing
	// 		Checkpoint = 8 + 32
	// 		AttestationData = 8 + 8 + 32 + 2 * Checkpoint
	// 		IndexedAttestation = 2048 + AttestationData + 96
	// 		AttesterSlashing = 2 * IndexedAttestation
	// 		# Attestation
	// 		Attestation = 2048 + AttestationData + 96
	// 		# Deposit
	// 		DepositData = 48 + 32 + 8 + 96
	// 		Deposit = 33 * 32 + DepositData
	// 		# SignedVoluntaryExit
	// 		VoluntaryExit = 8 + 8
	// 		SignedVoluntaryExit = VoluntaryExit + 96
	// 		# SyncAggregate
	// 		SyncAggregate = 64 + 96
	// 		# ExecutionPayload
	// 		Withdrawal = 8 + 8 + 20 + 8
	// 		ExecutionPayload = 32 + 20 + 32 + 32 + 256 + 32 + \
	// 							8 + 8 + 8 + 8 + 32 + 32 + 32 + \
	// 							MaxTransactions + 16 * Withdrawal + 8 + 8
	// 		# SignedBLSToExecutionChange
	// 		BLSToExecutionChange = 8 + 48 + 20
	// 		SignedBLSToExecutionChange = BLSToExecutionChange + 96
	// 		# KZGCommitment
	// 		KZGCommitment = 48
	// 		# BeaconBlodyBody
	// 		BeaconBlockBody = 96 + ETH1Data + Graffiti + 16 * ProposerSlashing + \
	// 						2 * AttesterSlashing + 128 * Attestation + \
	// 						16 * Deposit + 16 * SignedVoluntaryExit + \
	// 						SyncAggregate + ExecutionPayload + \
	// 						SignedBLSToExecutionChange + 4096
	// 		# BeaconBlock
	// 		BeaconBlock = 8 + 8 + 32 + 32 + BeaconBlockBody
	// 		# deneb.BlockContents
	// 		denebBlockContents = BeaconBlock + 6 * 48 + 6 * 131072
	DataSSZ []byte `ssz-max:"3620872"` // 2^22
}

func CreateConsensusData(rawSSZ []byte) (*ConsensusData, error) {
	cd := &ConsensusData{}
	err := cd.Decode(rawSSZ)
	return cd, err
}

func (cid *ConsensusData) Validate() error {
	switch cid.Duty.Type {
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

// GetBlockData ISSUE 221: GetBlockData/GetBlindedBlockData return versioned block only
func (ci *ConsensusData) GetBlockData() (blk *api.VersionedProposal, signingRoot ssz.HashRoot, err error) {
	switch ci.Version {
	case spec.DataVersionCapella:
		ret := &capella.BeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &api.VersionedProposal{Capella: ret, Version: ci.Version}, ret, nil
	case spec.DataVersionDeneb:
		ret := &apiv1deneb.BlockContents{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &api.VersionedProposal{Deneb: ret, Version: ci.Version}, ret.Block, nil
	default:
		return nil, nil, errors.Errorf("unknown block version %s", ci.Version.String())
	}
}

// GetBlindedBlockData ISSUE 221: GetBlockData/GetBlindedBlockData return versioned block only
func (ci *ConsensusData) GetBlindedBlockData() (*api.VersionedBlindedProposal, ssz.HashRoot, error) {
	switch ci.Version {
	case spec.DataVersionCapella:
		ret := &apiv1capella.BlindedBeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &api.VersionedBlindedProposal{Capella: ret, Version: ci.Version}, ret, nil
	case spec.DataVersionDeneb:
		ret := &apiv1deneb.BlindedBeaconBlock{}
		if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal ssz")
		}
		return &api.VersionedBlindedProposal{Deneb: ret, Version: ci.Version}, ret, nil
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
