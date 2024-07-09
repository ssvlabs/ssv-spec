package testingutils

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/api"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types"
)

const (

	//Deneb Fork Epoch
	ForkEpochPraterDeneb = 231680

	// ForkEpochPraterCapella Goerli taken from https://github.com/ethereum/execution-specs/blob/37a8f892341eb000e56e962a051a87e05a2e4443/network-upgrades/mainnet-upgrades/shanghai.md?plain=1#L18
	ForkEpochPraterCapella = 162304

	TestingDutyEpochCapella         = ForkEpochPraterCapella
	TestingDutySlotCapella          = ForkEpochPraterCapella * 32
	TestingDutySlotCapellaNextEpoch = TestingDutySlotCapella + 32
	TestingDutySlotCapellaInvalid   = TestingDutySlotCapella + 50

	TestingDutyEpochDeneb         = ForkEpochPraterDeneb
	TestingDutySlotDeneb          = ForkEpochPraterDeneb * 32
	TestingDutySlotDenebNextEpoch = TestingDutySlotDeneb + 32
	TestingDutySlotDenebInvalid   = TestingDutySlotDeneb + 50
)

// SupportedBlockVersions is a list of supported regular/blinded beacon block versions by spec.
var SupportedBlockVersions = []spec.DataVersion{spec.DataVersionCapella, spec.DataVersionDeneb}

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

	default:
		panic("unsupported version")
	}

	return ret
}

var TestingBlindedBeaconBlockV = func(version spec.DataVersion) *api.VersionedBlindedProposal {
	switch version {
	case spec.DataVersionCapella:
		return &api.VersionedBlindedProposal{
			Version: version,
			Capella: TestingBlindedBeaconBlockCapella,
		}
	case spec.DataVersionDeneb:
		return &api.VersionedBlindedProposal{
			Version: version,
			Deneb:   TestingBlindedBeaconBlockDeneb,
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
		if vBlk.Capella == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Capella.MarshalSSZ()
	case spec.DataVersionDeneb:
		if vBlk.Deneb == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Deneb.MarshalSSZ()

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
	default:
		panic("unsupported version")
	}
}

var TestingSignedBlindedBeaconBlockV = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	vBlk := TestingBlindedBeaconBlockV(version)

	switch version {
	case spec.DataVersionCapella:
		if vBlk.Capella == nil {
			panic("empty block")
		}
		return &apiv1capella.SignedBlindedBeaconBlock{
			Message:   vBlk.Capella,
			Signature: signBeaconObject(vBlk.Capella, types.DomainProposer, ks),
		}
	case spec.DataVersionDeneb:
		if vBlk.Deneb == nil {
			panic("empty block")
		}
		return &apiv1deneb.SignedBlindedBeaconBlock{
			Message:   vBlk.Deneb,
			Signature: signBeaconObject(vBlk.Deneb, types.DomainProposer, ks),
		}
	default:
		panic("unsupported version")
	}
}

var TestingDutyEpochV = func(version spec.DataVersion) phase0.Epoch {
	switch version {
	case spec.DataVersionCapella:
		return TestingDutyEpochCapella
	case spec.DataVersionDeneb:
		return TestingDutyEpochDeneb

	default:
		panic("unsupported version")
	}
}

var TestingDutySlotV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionCapella:
		return TestingDutySlotCapella
	case spec.DataVersionDeneb:
		return TestingDutySlotDeneb

	default:
		panic("unsupported version")
	}
}

var VersionBySlot = func(slot phase0.Slot) spec.DataVersion {
	if slot < ForkEpochPraterDeneb*32 {
		return spec.DataVersionCapella
	}
	return spec.DataVersionDeneb
}

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

	default:
		panic("unsupported version")
	}

	return duty
}

var TestingInvalidDutySlotV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionCapella:
		return TestingDutySlotCapellaInvalid
	case spec.DataVersionDeneb:
		return TestingDutySlotDenebInvalid

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
