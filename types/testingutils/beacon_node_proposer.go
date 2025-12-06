package testingutils

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/api"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"
	apiv1fulu "github.com/attestantio/go-eth2-client/api/v1/fulu"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Versioned (Blinded) Beacon Block
// ==================================================

// SupportedBlockVersions is a list of supported regular/blinded beacon block versions by spec.
var SupportedBlockVersions = []spec.DataVersion{spec.DataVersionCapella, spec.DataVersionDeneb, spec.DataVersionElectra, spec.DataVersionFulu}

var TestingBeaconBlockV = func(version spec.DataVersion) *api.VersionedProposal {
	switch version {
	case spec.DataVersionCapella:
		return &api.VersionedProposal{
			Version: version,
			Capella: TestingBeaconBlockCapella,
		}
	case spec.DataVersionDeneb:
		return &api.VersionedProposal{
			Version: version,
			Deneb:   TestingBlockContentsDeneb,
		}
	case spec.DataVersionElectra:
		return &api.VersionedProposal{
			Version: version,
			Electra: TestingBlockContentsElectra,
		}
	case spec.DataVersionFulu:
		return &api.VersionedProposal{
			Version: version,
			Fulu:    TestingBlockContentsFulu,
		}
	default:
		panic("unsupported version")
	}
}

var TestingBeaconBlockBytesV = func(version spec.DataVersion) []byte {
	var ret []byte
	vBlk := TestingBeaconBlockV(version)

	switch version {
	case spec.DataVersionCapella:
		if vBlk.Capella == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Capella.MarshalSSZ()
	case spec.DataVersionDeneb:
		if vBlk.Deneb == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Deneb.MarshalSSZ()
	case spec.DataVersionElectra:
		if vBlk.Electra == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Electra.MarshalSSZ()
	case spec.DataVersionFulu:
		if vBlk.Fulu == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Fulu.MarshalSSZ()
	default:
		panic("unsupported version")
	}

	return ret
}

var TestingBlindedBeaconBlockV = func(version spec.DataVersion) *api.VersionedProposal {
	switch version {
	case spec.DataVersionCapella:
		return &api.VersionedProposal{
			Version:        version,
			Blinded:        true,
			CapellaBlinded: TestingBlindedBeaconBlockCapella,
		}
	case spec.DataVersionDeneb:
		return &api.VersionedProposal{
			Version:      version,
			Blinded:      true,
			DenebBlinded: TestingBlindedBeaconBlockDeneb,
		}
	case spec.DataVersionElectra:
		return &api.VersionedProposal{
			Version:        version,
			Blinded:        true,
			ElectraBlinded: TestingBlindedBeaconBlockElectra,
		}
	case spec.DataVersionFulu:
		return &api.VersionedProposal{
			Version:     version,
			Blinded:     true,
			FuluBlinded: TestingBlindedBeaconBlockFulu,
		}
	default:
		panic("unsupported version")
	}
}

var TestingBlindedBeaconBlockBytesV = func(version spec.DataVersion) []byte {
	var ret []byte
	vBlk := TestingBlindedBeaconBlockV(version)

	switch version {
	case spec.DataVersionCapella:
		if vBlk.CapellaBlinded == nil {
			panic("empty block")
		}
		ret, _ = vBlk.CapellaBlinded.MarshalSSZ()
	case spec.DataVersionDeneb:
		if vBlk.DenebBlinded == nil {
			panic("empty block")
		}
		ret, _ = vBlk.DenebBlinded.MarshalSSZ()
	case spec.DataVersionElectra:
		if vBlk.ElectraBlinded == nil {
			panic("empty block")
		}
		ret, _ = vBlk.ElectraBlinded.MarshalSSZ()
	case spec.DataVersionFulu:
		if vBlk.FuluBlinded == nil {
			panic("empty block")
		}
		ret, _ = vBlk.FuluBlinded.MarshalSSZ()
	default:
		panic("unsupported version")
	}

	return ret
}

var TestingWrongBeaconBlockV = func(version spec.DataVersion) *api.VersionedProposal {
	blkByts := TestingBeaconBlockBytesV(version)

	switch version {
	case spec.DataVersionCapella:
		ret := &capella.BeaconBlock{}
		if err := ret.UnmarshalSSZ(blkByts); err != nil {
			panic(err.Error())
		}
		ret.Slot = TestingDutySlotCapella + 100
		return &api.VersionedProposal{
			Version: version,
			Capella: ret,
		}
	case spec.DataVersionDeneb:
		ret := &apiv1deneb.BlockContents{}
		if err := ret.UnmarshalSSZ(blkByts); err != nil {
			panic(err.Error())
		}
		ret.Block.Slot = TestingDutySlotDeneb + 100
		return &api.VersionedProposal{
			Version: version,
			Deneb:   ret,
		}
	case spec.DataVersionElectra:
		ret := &apiv1electra.BlockContents{}
		if err := ret.UnmarshalSSZ(blkByts); err != nil {
			panic(err.Error())
		}
		ret.Block.Slot = TestingDutySlotElectra + 100
		return &api.VersionedProposal{
			Version: version,
			Electra: ret,
		}
	case spec.DataVersionFulu:
		ret := &apiv1fulu.BlockContents{}
		if err := ret.UnmarshalSSZ(blkByts); err != nil {
			panic(err.Error())
		}
		ret.Block.Slot = TestingDutySlotFulu + 100
		return &api.VersionedProposal{
			Version: version,
			Fulu:    ret,
		}
	default:
		panic("unsupported version")
	}
}

var TestingSignedBeaconBlockV = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	vBlk := TestingBeaconBlockV(version)

	switch version {
	case spec.DataVersionCapella:
		if vBlk.Capella == nil {
			panic("empty block")
		}
		return &capella.SignedBeaconBlock{
			Message:   vBlk.Capella,
			Signature: signBeaconObject(vBlk.Capella, types.DomainProposer, ks),
		}
	case spec.DataVersionDeneb:
		if vBlk.Deneb == nil {
			panic("empty block contents")
		}
		if vBlk.Deneb.Block == nil {
			panic("empty block")
		}
		return &apiv1deneb.SignedBlockContents{
			SignedBlock: &deneb.SignedBeaconBlock{
				Message:   vBlk.Deneb.Block,
				Signature: signBeaconObject(vBlk.Deneb.Block, types.DomainProposer, ks),
			},
			KZGProofs: vBlk.Deneb.KZGProofs,
			Blobs:     vBlk.Deneb.Blobs,
		}
	case spec.DataVersionElectra:
		if vBlk.Electra == nil {
			panic("empty block contents")
		}
		if vBlk.Electra.Block == nil {
			panic("empty block")
		}
		return &apiv1electra.SignedBlockContents{
			SignedBlock: &electra.SignedBeaconBlock{
				Message:   vBlk.Electra.Block,
				Signature: signBeaconObject(vBlk.Electra.Block, types.DomainProposer, ks),
			},
			KZGProofs: vBlk.Electra.KZGProofs,
			Blobs:     vBlk.Electra.Blobs,
		}
	case spec.DataVersionFulu:
		if vBlk.Fulu == nil {
			panic("empty block contents")
		}
		if vBlk.Fulu.Block == nil {
			panic("empty block")
		}
		return &apiv1fulu.SignedBlockContents{
			SignedBlock: &electra.SignedBeaconBlock{
				Message:   vBlk.Fulu.Block,
				Signature: signBeaconObject(vBlk.Fulu.Block, types.DomainProposer, ks),
			},
			KZGProofs: vBlk.Fulu.KZGProofs,
			Blobs:     vBlk.Fulu.Blobs,
		}
	default:
		panic("unsupported version")
	}
}

var TestingSignedBlindedBeaconBlockV = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	vBlk := TestingBlindedBeaconBlockV(version)

	switch version {
	case spec.DataVersionCapella:
		if vBlk.CapellaBlinded == nil {
			panic("empty block")
		}
		return &apiv1capella.SignedBlindedBeaconBlock{
			Message:   vBlk.CapellaBlinded,
			Signature: signBeaconObject(vBlk.CapellaBlinded, types.DomainProposer, ks),
		}
	case spec.DataVersionDeneb:
		if vBlk.DenebBlinded == nil {
			panic("empty block")
		}
		return &apiv1deneb.SignedBlindedBeaconBlock{
			Message:   vBlk.DenebBlinded,
			Signature: signBeaconObject(vBlk.DenebBlinded, types.DomainProposer, ks),
		}
	case spec.DataVersionElectra:
		if vBlk.ElectraBlinded == nil {
			panic("empty block")
		}
		return &apiv1electra.SignedBlindedBeaconBlock{
			Message:   vBlk.ElectraBlinded,
			Signature: signBeaconObject(vBlk.ElectraBlinded, types.DomainProposer, ks),
		}
	case spec.DataVersionFulu:
		if vBlk.FuluBlinded == nil {
			panic("empty block")
		}
		return &apiv1electra.SignedBlindedBeaconBlock{
			Message:   vBlk.FuluBlinded,
			Signature: signBeaconObject(vBlk.FuluBlinded, types.DomainProposer, ks),
		}
	default:
		panic("unsupported version")
	}
}

var TestingBeaconBlockCapella = func() *capella.BeaconBlock {
	var res capella.BeaconBlock
	err := json.Unmarshal(capellaBlock, &res)
	if err != nil {
		panic(err)
	}
	// using TestingDutySlotCapella to keep the consistency with TestingProposerDutyV Capella slot
	res.Slot = TestingDutySlotCapella
	return &res
}()

var TestingBlindedBeaconBlockCapella = func() *apiv1capella.BlindedBeaconBlock {
	fullBlk := TestingBeaconBlockCapella
	txRoot, _ := types.SSZTransactions(fullBlk.Body.ExecutionPayload.Transactions).HashTreeRoot()
	withdrawalsRoot, _ := types.SSZWithdrawals(fullBlk.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
	ret := &apiv1capella.BlindedBeaconBlock{
		Slot:          fullBlk.Slot,
		ProposerIndex: fullBlk.ProposerIndex,
		ParentRoot:    fullBlk.ParentRoot,
		StateRoot:     fullBlk.StateRoot,
		Body: &apiv1capella.BlindedBeaconBlockBody{
			RANDAOReveal:      fullBlk.Body.RANDAOReveal,
			ETH1Data:          fullBlk.Body.ETH1Data,
			Graffiti:          fullBlk.Body.Graffiti,
			ProposerSlashings: fullBlk.Body.ProposerSlashings,
			AttesterSlashings: fullBlk.Body.AttesterSlashings,
			Attestations:      fullBlk.Body.Attestations,
			Deposits:          fullBlk.Body.Deposits,
			VoluntaryExits:    fullBlk.Body.VoluntaryExits,
			SyncAggregate:     fullBlk.Body.SyncAggregate,
			ExecutionPayloadHeader: &capella.ExecutionPayloadHeader{
				ParentHash:       fullBlk.Body.ExecutionPayload.ParentHash,
				FeeRecipient:     fullBlk.Body.ExecutionPayload.FeeRecipient,
				StateRoot:        fullBlk.Body.ExecutionPayload.StateRoot,
				ReceiptsRoot:     fullBlk.Body.ExecutionPayload.ReceiptsRoot,
				LogsBloom:        fullBlk.Body.ExecutionPayload.LogsBloom,
				PrevRandao:       fullBlk.Body.ExecutionPayload.PrevRandao,
				BlockNumber:      fullBlk.Body.ExecutionPayload.BlockNumber,
				GasLimit:         fullBlk.Body.ExecutionPayload.GasLimit,
				GasUsed:          fullBlk.Body.ExecutionPayload.GasUsed,
				Timestamp:        fullBlk.Body.ExecutionPayload.Timestamp,
				ExtraData:        fullBlk.Body.ExecutionPayload.ExtraData,
				BaseFeePerGas:    fullBlk.Body.ExecutionPayload.BaseFeePerGas,
				BlockHash:        fullBlk.Body.ExecutionPayload.BlockHash,
				TransactionsRoot: txRoot,
				WithdrawalsRoot:  withdrawalsRoot,
			},
			BLSToExecutionChanges: fullBlk.Body.BLSToExecutionChanges,
		},
	}

	return ret
}()

var TestingBlockContentsDeneb = func() *apiv1deneb.BlockContents {
	var res apiv1deneb.BlockContents
	if err := json.Unmarshal(denebBlockContents, &res); err != nil {
		panic(err)
	}
	// using ForkEpochPraterDeneb to keep the consistency with TestingProposerDutyV Deneb slot
	res.Block.Slot = ForkEpochPraterDeneb
	return &res
}()

var TestingBlindedBeaconBlockDeneb = func() *apiv1deneb.BlindedBeaconBlock {
	blockContents := TestingBlockContentsDeneb
	txRoot, _ := types.SSZTransactions(blockContents.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
	withdrawalsRoot, _ := types.SSZWithdrawals(blockContents.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
	ret := &apiv1deneb.BlindedBeaconBlock{
		Slot:          blockContents.Block.Slot,
		ProposerIndex: blockContents.Block.ProposerIndex,
		ParentRoot:    blockContents.Block.ParentRoot,
		StateRoot:     blockContents.Block.StateRoot,
		Body: &apiv1deneb.BlindedBeaconBlockBody{
			RANDAOReveal:      blockContents.Block.Body.RANDAOReveal,
			ETH1Data:          blockContents.Block.Body.ETH1Data,
			Graffiti:          blockContents.Block.Body.Graffiti,
			ProposerSlashings: blockContents.Block.Body.ProposerSlashings,
			AttesterSlashings: blockContents.Block.Body.AttesterSlashings,
			Attestations:      blockContents.Block.Body.Attestations,
			Deposits:          blockContents.Block.Body.Deposits,
			VoluntaryExits:    blockContents.Block.Body.VoluntaryExits,
			SyncAggregate:     blockContents.Block.Body.SyncAggregate,
			ExecutionPayloadHeader: &deneb.ExecutionPayloadHeader{
				ParentHash:       blockContents.Block.Body.ExecutionPayload.ParentHash,
				FeeRecipient:     blockContents.Block.Body.ExecutionPayload.FeeRecipient,
				StateRoot:        blockContents.Block.Body.ExecutionPayload.StateRoot,
				ReceiptsRoot:     blockContents.Block.Body.ExecutionPayload.ReceiptsRoot,
				LogsBloom:        blockContents.Block.Body.ExecutionPayload.LogsBloom,
				PrevRandao:       blockContents.Block.Body.ExecutionPayload.PrevRandao,
				BlockNumber:      blockContents.Block.Body.ExecutionPayload.BlockNumber,
				GasLimit:         blockContents.Block.Body.ExecutionPayload.GasLimit,
				GasUsed:          blockContents.Block.Body.ExecutionPayload.GasUsed,
				Timestamp:        blockContents.Block.Body.ExecutionPayload.Timestamp,
				ExtraData:        blockContents.Block.Body.ExecutionPayload.ExtraData,
				BaseFeePerGas:    blockContents.Block.Body.ExecutionPayload.BaseFeePerGas,
				BlockHash:        blockContents.Block.Body.ExecutionPayload.BlockHash,
				TransactionsRoot: txRoot,
				WithdrawalsRoot:  withdrawalsRoot,
				BlobGasUsed:      blockContents.Block.Body.ExecutionPayload.BlobGasUsed,
				ExcessBlobGas:    blockContents.Block.Body.ExecutionPayload.ExcessBlobGas,
			},
			BLSToExecutionChanges: blockContents.Block.Body.BLSToExecutionChanges,
			BlobKZGCommitments:    blockContents.Block.Body.BlobKZGCommitments,
		},
	}

	return ret
}()

var TestingBlockContentsElectra = func() *apiv1electra.BlockContents {
	var beaconBlock *electra.BeaconBlock
	if err := json.Unmarshal(electraBeaconBlock, &beaconBlock); err != nil {
		panic(err)
	}

	// using ForkEpochPraterElectra to keep the consistency with TestingProposerDutyV Electra slot
	beaconBlock.Slot = ForkEpochPraterElectra

	denebBlockContents := TestingBlockContentsDeneb

	return &apiv1electra.BlockContents{
		Block:     beaconBlock,
		KZGProofs: denebBlockContents.KZGProofs,
		Blobs:     denebBlockContents.Blobs,
	}
}()

var TestingBlindedBeaconBlockElectra = func() *apiv1electra.BlindedBeaconBlock {
	blockContents := TestingBlockContentsElectra
	txRoot, _ := types.SSZTransactions(blockContents.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
	withdrawalsRoot, _ := types.SSZWithdrawals(blockContents.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
	ret := &apiv1electra.BlindedBeaconBlock{
		Slot:          blockContents.Block.Slot,
		ProposerIndex: blockContents.Block.ProposerIndex,
		ParentRoot:    blockContents.Block.ParentRoot,
		StateRoot:     blockContents.Block.StateRoot,
		Body: &apiv1electra.BlindedBeaconBlockBody{
			RANDAOReveal:      blockContents.Block.Body.RANDAOReveal,
			ETH1Data:          blockContents.Block.Body.ETH1Data,
			Graffiti:          blockContents.Block.Body.Graffiti,
			ProposerSlashings: blockContents.Block.Body.ProposerSlashings,
			AttesterSlashings: blockContents.Block.Body.AttesterSlashings,
			Attestations:      blockContents.Block.Body.Attestations,
			Deposits:          blockContents.Block.Body.Deposits,
			VoluntaryExits:    blockContents.Block.Body.VoluntaryExits,
			SyncAggregate:     blockContents.Block.Body.SyncAggregate,
			ExecutionPayloadHeader: &deneb.ExecutionPayloadHeader{
				ParentHash:       blockContents.Block.Body.ExecutionPayload.ParentHash,
				FeeRecipient:     blockContents.Block.Body.ExecutionPayload.FeeRecipient,
				StateRoot:        blockContents.Block.Body.ExecutionPayload.StateRoot,
				ReceiptsRoot:     blockContents.Block.Body.ExecutionPayload.ReceiptsRoot,
				LogsBloom:        blockContents.Block.Body.ExecutionPayload.LogsBloom,
				PrevRandao:       blockContents.Block.Body.ExecutionPayload.PrevRandao,
				BlockNumber:      blockContents.Block.Body.ExecutionPayload.BlockNumber,
				GasLimit:         blockContents.Block.Body.ExecutionPayload.GasLimit,
				GasUsed:          blockContents.Block.Body.ExecutionPayload.GasUsed,
				Timestamp:        blockContents.Block.Body.ExecutionPayload.Timestamp,
				ExtraData:        blockContents.Block.Body.ExecutionPayload.ExtraData,
				BaseFeePerGas:    blockContents.Block.Body.ExecutionPayload.BaseFeePerGas,
				BlockHash:        blockContents.Block.Body.ExecutionPayload.BlockHash,
				TransactionsRoot: txRoot,
				WithdrawalsRoot:  withdrawalsRoot,
				BlobGasUsed:      blockContents.Block.Body.ExecutionPayload.BlobGasUsed,
				ExcessBlobGas:    blockContents.Block.Body.ExecutionPayload.ExcessBlobGas,
			},
			BLSToExecutionChanges: blockContents.Block.Body.BLSToExecutionChanges,
			BlobKZGCommitments:    blockContents.Block.Body.BlobKZGCommitments,
			ExecutionRequests:     blockContents.Block.Body.ExecutionRequests,
		},
	}

	return ret
}()

var TestingBlockContentsFulu = func() *apiv1fulu.BlockContents {
	var beaconBlock *electra.BeaconBlock
	if err := json.Unmarshal(electraBeaconBlock, &beaconBlock); err != nil {
		panic(err)
	}

	// using ForkEpochFulu to keep the consistency with TestingProposerDutyV Fulu slot
	beaconBlock.Slot = ForkEpochFulu

	denebBlockContents := TestingBlockContentsDeneb

	return &apiv1fulu.BlockContents{
		Block:     beaconBlock,
		KZGProofs: denebBlockContents.KZGProofs,
		Blobs:     denebBlockContents.Blobs,
	}
}()

var TestingBlindedBeaconBlockFulu = func() *apiv1electra.BlindedBeaconBlock {
	blockContents := TestingBlockContentsFulu
	txRoot, _ := types.SSZTransactions(blockContents.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
	withdrawalsRoot, _ := types.SSZWithdrawals(blockContents.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
	ret := &apiv1electra.BlindedBeaconBlock{
		Slot:          blockContents.Block.Slot,
		ProposerIndex: blockContents.Block.ProposerIndex,
		ParentRoot:    blockContents.Block.ParentRoot,
		StateRoot:     blockContents.Block.StateRoot,
		Body: &apiv1electra.BlindedBeaconBlockBody{
			RANDAOReveal:      blockContents.Block.Body.RANDAOReveal,
			ETH1Data:          blockContents.Block.Body.ETH1Data,
			Graffiti:          blockContents.Block.Body.Graffiti,
			ProposerSlashings: blockContents.Block.Body.ProposerSlashings,
			AttesterSlashings: blockContents.Block.Body.AttesterSlashings,
			Attestations:      blockContents.Block.Body.Attestations,
			Deposits:          blockContents.Block.Body.Deposits,
			VoluntaryExits:    blockContents.Block.Body.VoluntaryExits,
			SyncAggregate:     blockContents.Block.Body.SyncAggregate,
			ExecutionPayloadHeader: &deneb.ExecutionPayloadHeader{
				ParentHash:       blockContents.Block.Body.ExecutionPayload.ParentHash,
				FeeRecipient:     blockContents.Block.Body.ExecutionPayload.FeeRecipient,
				StateRoot:        blockContents.Block.Body.ExecutionPayload.StateRoot,
				ReceiptsRoot:     blockContents.Block.Body.ExecutionPayload.ReceiptsRoot,
				LogsBloom:        blockContents.Block.Body.ExecutionPayload.LogsBloom,
				PrevRandao:       blockContents.Block.Body.ExecutionPayload.PrevRandao,
				BlockNumber:      blockContents.Block.Body.ExecutionPayload.BlockNumber,
				GasLimit:         blockContents.Block.Body.ExecutionPayload.GasLimit,
				GasUsed:          blockContents.Block.Body.ExecutionPayload.GasUsed,
				Timestamp:        blockContents.Block.Body.ExecutionPayload.Timestamp,
				ExtraData:        blockContents.Block.Body.ExecutionPayload.ExtraData,
				BaseFeePerGas:    blockContents.Block.Body.ExecutionPayload.BaseFeePerGas,
				BlockHash:        blockContents.Block.Body.ExecutionPayload.BlockHash,
				TransactionsRoot: txRoot,
				WithdrawalsRoot:  withdrawalsRoot,
				BlobGasUsed:      blockContents.Block.Body.ExecutionPayload.BlobGasUsed,
				ExcessBlobGas:    blockContents.Block.Body.ExecutionPayload.ExcessBlobGas,
			},
			BLSToExecutionChanges: blockContents.Block.Body.BLSToExecutionChanges,
			BlobKZGCommitments:    blockContents.Block.Body.BlobKZGCommitments,
			ExecutionRequests:     blockContents.Block.Body.ExecutionRequests,
		},
	}

	return ret
}()

// ==================================================
// Versioned Proposer Duty
// ==================================================

var TestingProposerDutyV = func(version spec.DataVersion) *types.ValidatorDuty {
	duty := &types.ValidatorDuty{
		Type:           types.BNRoleProposer,
		PubKey:         TestingValidatorPubKey,
		Slot:           TestingDutySlotV(version),
		ValidatorIndex: TestingValidatorIndex,
		// ISSUE 233: We are initializing unused struct fields here
		CommitteeIndex:          3,
		CommitteesAtSlot:        36,
		CommitteeLength:         128,
		ValidatorCommitteeIndex: 11,
	}

	return duty
}

var TestingProposerDutyNextEpochV = func(version spec.DataVersion) *types.ValidatorDuty {
	duty := &types.ValidatorDuty{
		Type:           types.BNRoleProposer,
		PubKey:         TestingValidatorPubKey,
		ValidatorIndex: TestingValidatorIndex,
		// ISSUE 233: We are initializing unused struct fields here
		CommitteeIndex:          3,
		CommitteesAtSlot:        36,
		CommitteeLength:         128,
		ValidatorCommitteeIndex: 11,
	}

	switch version {
	case spec.DataVersionCapella:
		duty.Slot = TestingDutySlotCapellaNextEpoch
	case spec.DataVersionDeneb:
		duty.Slot = TestingDutySlotDenebNextEpoch
	case spec.DataVersionElectra:
		duty.Slot = TestingDutySlotElectraNextEpoch
	case spec.DataVersionFulu:
		duty.Slot = TestingDutySlotFuluNextEpoch

	default:
		panic("unsupported version")
	}

	return duty
}

var TestingProposerDutyFirstSlot = types.ValidatorDuty{
	Type:           types.BNRoleProposer,
	PubKey:         TestingValidatorPubKey,
	Slot:           0,
	ValidatorIndex: TestingValidatorIndex,
}
