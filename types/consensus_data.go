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

// BeaconVote is used as the data to be agreed on consensus for the CommitteeRunner
type BeaconVote struct {
	BlockRoot phase0.Root `ssz-size:"32"`
	Source    *phase0.Checkpoint
	Target    *phase0.Checkpoint
}

// Encode the BeaconVote object
func (b *BeaconVote) Encode() ([]byte, error) {
	return b.MarshalSSZ()
}

// Decode the BeaconVote object
func (b *BeaconVote) Decode(data []byte) error {
	return b.UnmarshalSSZ(data)
}

// ValidatorConsensusData holds all relevant duty and data Decided on by consensus
type ValidatorConsensusData struct {
	// Duty max size is
	// 			8 + 48 + 6*8 + 13*8 + 1 = 209
	Duty    ValidatorDuty
	Version spec.DataVersion
	// DataSSZ's max size if the size of the biggest object Deneb.BlockContents.
	// Per definition, Deneb.BlockContents has a field for transaction of size 2^50.
	// We do not need to support such a big DataSSZ size as 2^50 represents 1000X the actual block gas limit
	// Upcoming 40M gas limit produces 40M / 16 (call data cost) = 2,500,000 bytes (https://eips.ethereum.org/EIPS/eip-4488)
	// Explanation on why transaction sizes are so big https://github.com/ethereum/consensus-specs/pull/2686
	// Adding to the rest of the data (see script below), we have: 1,315,964 + 2,500,000  = 3,815,964 bytes ~<= 2^22
	// Python script for Deneb.BlockContents without transactions:
	// 		# Constants
	// 		KZG_PROOFS_SIZE = 6 * 48  # KZGProofs size
	// 		BLOBS_SIZE = 6 * 131072  # Blobs size
	// 		BEACON_BLOCK_OVERHEAD = 2 * 32 + 2 * 8  # Additional overhead for BeaconBlock
	// 		# Components of BeaconBlockBody
	// 		ETH1_DATA_SIZE = 96 + 2 * 32 + 8 + 32  # ETH1Data
	// 		PROPOSER_SLASHING_SIZE = 16 * (2 * (96 + 3 * 32 + 2 * 8))  # ProposerSlashing
	// 		ATTESTER_SLASHING_SIZE = 2 * (2 * (2048 + 96 + (2 * 8 + 32 + 2 * (8 + 32))))  # AttesterSlashing
	// 		ATTESTATION_SIZE = 128 * (2048 + 96 + (2 * 8 + 32 + 2 * (8 + 32)))  # Attestation
	// 		DEPOSIT_SIZE = 16 * (33 * 32 + 48 + 32 + 8 + 96)  # Deposit
	// 		SIGNED_VOLUNTARY_EXIT_SIZE = 16 * (96 + 2 * 8)  # SignedVoluntaryExit
	// 		SYNC_AGGREGATE_SIZE = 64 + 96  # SyncAggregate
	// 		EXECUTION_PAYLOAD_NO_TRANSACTIONS = 32 + 20 + 2*32 + 256 + 32 + 4*8 + 3*32 + 16 * (2*8 + 20 + 8) + 8 + 8
	// 		BLS_TO_EXECUTION_CHANGES_SIZE = 16 * (96 + (8 + 48 + 20))  # BLSToExecutionChanges
	// 		KZG_COMMITMENT_SIZE = 4096 * 48  # KZGCommitment
	// 		# BeaconBlockBody total size without transactions
	// 		beacon_block_body_size_without_transactions = (
	// 		    ETH1_DATA_SIZE + PROPOSER_SLASHING_SIZE + ATTESTER_SLASHING_SIZE +
	// 		    ATTESTATION_SIZE + DEPOSIT_SIZE + SIGNED_VOLUNTARY_EXIT_SIZE +
	// 		    SYNC_AGGREGATE_SIZE + EXECUTION_PAYLOAD_NO_TRANSACTIONS + BLS_TO_EXECUTION_CHANGES_SIZE + KZG_COMMITMENT_SIZE
	// 		)
	// 		# Total size of Deneb.BlockContents and BeaconBlock without transactions
	// 		total_size_without_execution_payload = KZG_PROOFS_SIZE + BLOBS_SIZE + BEACON_BLOCK_OVERHEAD + beacon_block_body_size_without_transactions
	DataSSZ []byte `ssz-max:"4194304"` // 2^22
}

func (cid *ValidatorConsensusData) Validate() error {
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
func (ci *ValidatorConsensusData) GetBlockData() (blk *api.VersionedProposal, signingRoot ssz.HashRoot, err error) {
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
func (ci *ValidatorConsensusData) GetBlindedBlockData() (*api.VersionedBlindedProposal, ssz.HashRoot, error) {
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

func (ci *ValidatorConsensusData) GetAggregateAndProof() (*phase0.AggregateAndProof, error) {
	ret := &phase0.AggregateAndProof{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ssz")
	}
	return ret, nil
}

func (ci *ValidatorConsensusData) GetSyncCommitteeContributions() (Contributions, error) {
	ret := Contributions{}
	if err := ret.UnmarshalSSZ(ci.DataSSZ); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ssz")
	}
	return ret, nil
}

func (cid *ValidatorConsensusData) Encode() ([]byte, error) {
	return cid.MarshalSSZ()
}

func (cid *ValidatorConsensusData) Decode(data []byte) error {
	return cid.UnmarshalSSZ(data)
}
