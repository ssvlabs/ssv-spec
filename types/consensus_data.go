package types

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"
	apiv1fulu "github.com/attestantio/go-eth2-client/api/v1/fulu"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
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
	// Adding to the rest of the data (see script below), we have: 3,291,849 + 2,500,000  = 5,791,849 bytes ~<= 2^23
	// Python script for Deneb.BlockContents without transactions:
	// 		# Constants
	// 		KZG_PROOFS_SIZE = 9 * 48  # KZGProofs size
	// 		BLOBS_SIZE = 9 * 131072  # Blobs size
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
	//		EXECUTION_REQUESTS_SIZE = (1 + 8192 * (1 + 48 + 32 + 8 + 96 + 8)) + (1 + 16 * (1 + 20 + 48 + 8)) + (1 + 2 * (1 + 20 +  48 + 48)) # Deposits + Withdrawls + Consolidations
	// 		# BeaconBlockBody total size without transactions
	// 		beacon_block_body_size_without_transactions = (
	// 		    ETH1_DATA_SIZE + PROPOSER_SLASHING_SIZE + ATTESTER_SLASHING_SIZE +
	// 		    ATTESTATION_SIZE + DEPOSIT_SIZE + SIGNED_VOLUNTARY_EXIT_SIZE +
	// 		    SYNC_AGGREGATE_SIZE + EXECUTION_PAYLOAD_NO_TRANSACTIONS + BLS_TO_EXECUTION_CHANGES_SIZE + KZG_COMMITMENT_SIZE + EXECUTION_REQUESTS_SIZE
	// 		)
	// 		# Total size of Deneb.BlockContents and BeaconBlock without transactions
	// 		total_size_without_execution_payload = KZG_PROOFS_SIZE + BLOBS_SIZE + BEACON_BLOCK_OVERHEAD + beacon_block_body_size_without_transactions
	//		print(total_size_without_execution_payload)
	DataSSZ []byte `ssz-max:"8388608"` // 2^23 to account for potential gas limit increases
}

func (cd *ValidatorConsensusData) Validate() error {
	switch cd.Duty.Type {
	case BNRoleAggregator:
		if _, _, err := cd.GetAggregateAndProof(); err != nil {
			return err
		}
	case BNRoleProposer:
		if _, _, err := cd.GetBlockData(); err != nil {
			return err
		}
	case BNRoleSyncCommitteeContribution:
		if _, err := cd.GetSyncCommitteeContributions(); err != nil {
			return err
		}
	case BNRoleValidatorRegistration:
		return NewError(ValidatorRegistrationNoConsensusDataErrorCode, "validator registration has no consensus data")
	case BNRoleVoluntaryExit:
		return NewError(ValidatorExitNoConsensusDataErrorCode, "voluntary exit has no consensus data")
	default:
		return NewError(UnknownDutyRoleDataErrorCode, "unknown duty role")
	}
	return nil
}

// GetBlockData returns block data for both blinded and regular blocks
func (cd *ValidatorConsensusData) GetBlockData() (blk *api.VersionedProposal, signingRoot ssz.HashRoot, err error) {
	switch cd.Version {
	case spec.DataVersionCapella:
		blindedBlock := &apiv1capella.BlindedBeaconBlock{}
		blindedErr := blindedBlock.UnmarshalSSZ(cd.DataSSZ)
		if blindedErr == nil {
			return &api.VersionedProposal{Version: cd.Version, Blinded: true, CapellaBlinded: blindedBlock}, blindedBlock, nil
		}

		regularBlock := &capella.BeaconBlock{}
		regularErr := regularBlock.UnmarshalSSZ(cd.DataSSZ)
		if regularErr == nil {
			return &api.VersionedProposal{Capella: regularBlock, Version: cd.Version}, regularBlock, nil
		}

		return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz (blinded err: %w, regular err: %w)", blindedErr, regularErr))
	case spec.DataVersionDeneb:
		blindedBlock := &apiv1deneb.BlindedBeaconBlock{}
		blindedErr := blindedBlock.UnmarshalSSZ(cd.DataSSZ)
		if blindedErr == nil {
			return &api.VersionedProposal{Version: cd.Version, Blinded: true, DenebBlinded: blindedBlock}, blindedBlock, nil
		}

		regularContents := &apiv1deneb.BlockContents{}
		regularErr := regularContents.UnmarshalSSZ(cd.DataSSZ)
		if regularErr == nil {
			return &api.VersionedProposal{Deneb: regularContents, Version: cd.Version}, regularContents.Block, nil
		}

		return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz (blinded err: %w, regular err: %w)", blindedErr, regularErr))
	case spec.DataVersionElectra:
		blindedBlock := &apiv1electra.BlindedBeaconBlock{}
		blindedErr := blindedBlock.UnmarshalSSZ(cd.DataSSZ)
		if blindedErr == nil {
			return &api.VersionedProposal{Version: cd.Version, Blinded: true, ElectraBlinded: blindedBlock}, blindedBlock, nil
		}

		regularContents := &apiv1electra.BlockContents{}
		regularErr := regularContents.UnmarshalSSZ(cd.DataSSZ)
		if regularErr == nil {
			return &api.VersionedProposal{Electra: regularContents, Version: cd.Version}, regularContents.Block, nil
		}

		return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz (blinded err: %w, regular err: %w)", blindedErr, regularErr))
	case spec.DataVersionFulu:
		blindedBlock := &apiv1electra.BlindedBeaconBlock{}
		blindedErr := blindedBlock.UnmarshalSSZ(cd.DataSSZ)
		if blindedErr == nil {
			return &api.VersionedProposal{Version: cd.Version, Blinded: true, FuluBlinded: blindedBlock}, blindedBlock, nil
		}

		regularContents := &apiv1fulu.BlockContents{}
		regularErr := regularContents.UnmarshalSSZ(cd.DataSSZ)
		if regularErr == nil {
			return &api.VersionedProposal{Fulu: regularContents, Version: cd.Version}, regularContents.Block, nil
		}

		return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz (blinded err: %w, regular err: %w)", blindedErr, regularErr))
	default:
		return nil, nil, WrapError(UnknownBlockVersionErrorCode, fmt.Errorf("unknown block version %d", cd.Version))
	}
}

func (cd *ValidatorConsensusData) GetAggregateAndProof() (*spec.VersionedAggregateAndProof, ssz.HashRoot, error) {
	switch cd.Version {
	case spec.DataVersionPhase0:
		ret := &phase0.AggregateAndProof{}
		if err := ret.UnmarshalSSZ(cd.DataSSZ); err != nil {
			return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz: %w", err))
		}

		return &spec.VersionedAggregateAndProof{Version: cd.Version, Phase0: ret}, ret, nil
	case spec.DataVersionAltair:
		ret := &phase0.AggregateAndProof{}
		if err := ret.UnmarshalSSZ(cd.DataSSZ); err != nil {
			return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz: %w", err))
		}

		return &spec.VersionedAggregateAndProof{Version: cd.Version, Altair: ret}, ret, nil
	case spec.DataVersionBellatrix:
		ret := &phase0.AggregateAndProof{}
		if err := ret.UnmarshalSSZ(cd.DataSSZ); err != nil {
			return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz: %w", err))
		}

		return &spec.VersionedAggregateAndProof{Version: cd.Version, Bellatrix: ret}, ret, nil
	case spec.DataVersionCapella:
		ret := &phase0.AggregateAndProof{}
		if err := ret.UnmarshalSSZ(cd.DataSSZ); err != nil {
			return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz: %w", err))
		}

		return &spec.VersionedAggregateAndProof{Version: cd.Version, Capella: ret}, ret, nil
	case spec.DataVersionDeneb:
		ret := &phase0.AggregateAndProof{}
		if err := ret.UnmarshalSSZ(cd.DataSSZ); err != nil {
			return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz: %w", err))
		}

		return &spec.VersionedAggregateAndProof{Version: cd.Version, Deneb: ret}, ret, nil
	case spec.DataVersionElectra:
		ret := &electra.AggregateAndProof{}
		if err := ret.UnmarshalSSZ(cd.DataSSZ); err != nil {
			return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz: %w", err))
		}

		return &spec.VersionedAggregateAndProof{Version: cd.Version, Electra: ret}, ret, nil
	case spec.DataVersionFulu:
		ret := &electra.AggregateAndProof{}
		if err := ret.UnmarshalSSZ(cd.DataSSZ); err != nil {
			return nil, nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz: %w", err))
		}

		return &spec.VersionedAggregateAndProof{Version: cd.Version, Fulu: ret}, ret, nil
	default:
		return nil, nil, fmt.Errorf("unknown aggregate and proof version %d", cd.Version)
	}
}

func (cd *ValidatorConsensusData) GetSyncCommitteeContributions() (Contributions, error) {
	ret := Contributions{}
	if err := ret.UnmarshalSSZ(cd.DataSSZ); err != nil {
		return nil, WrapError(UnmarshalSSZErrorCode, fmt.Errorf("could not unmarshal ssz: %w", err))
	}
	return ret, nil
}

func (cd *ValidatorConsensusData) Encode() ([]byte, error) {
	return cd.MarshalSSZ()
}

func (cd *ValidatorConsensusData) Decode(data []byte) error {
	return cd.UnmarshalSSZ(data)
}
