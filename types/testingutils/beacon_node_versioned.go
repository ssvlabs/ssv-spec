package testingutils

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/api"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/bloxapp/ssv-spec/types"
)

const (
	// ForkEpochPraterCapella Goerli taken from https://github.com/ethereum/execution-specs/blob/37a8f892341eb000e56e962a051a87e05a2e4443/network-upgrades/mainnet-upgrades/shanghai.md?plain=1#L18
	ForkEpochPraterCapella = 162304

	// TestingDutySlotBellatrix keeping this value to not break the test roots
	TestingDutySlotBellatrix          = 12
	TestingDutySlotBellatrixNextEpoch = 50
	TestingDutySlotBellatrixInvalid   = 50
	TestingDutyEpochBellatrix         = 0

	TestingDutyEpochCapella         = ForkEpochPraterCapella
	TestingDutySlotCapella          = ForkEpochPraterCapella * 32
	TestingDutySlotCapellaNextEpoch = TestingDutySlotCapella + 32
	TestingDutySlotCapellaInvalid   = TestingDutySlotCapella + 50
)

// SupportedBlockVersions is a list of supported regular/blinded beacon block versions by spec.
var SupportedBlockVersions = []spec.DataVersion{spec.DataVersionBellatrix, spec.DataVersionCapella}

var TestingBeaconBlockV = func(version spec.DataVersion) *spec.VersionedBeaconBlock {
	switch version {
	case spec.DataVersionBellatrix:
		return &spec.VersionedBeaconBlock{
			Version:   version,
			Bellatrix: TestingBeaconBlock,
		}
	case spec.DataVersionCapella:
		return &spec.VersionedBeaconBlock{
			Version: version,
			Capella: TestingBeaconBlockCapella,
		}
	default:
		panic("unsupported version")
	}
}

var TestingBeaconBlockBytesV = func(version spec.DataVersion) []byte {
	var ret []byte
	vBlk := TestingBeaconBlockV(version)

	switch version {
	case spec.DataVersionBellatrix:
		if vBlk.Bellatrix == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Bellatrix.MarshalSSZ()
	case spec.DataVersionCapella:
		if vBlk.Capella == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Capella.MarshalSSZ()

	default:
		panic("unsupported version")
	}

	return ret
}

var TestingBlindedBeaconBlockV = func(version spec.DataVersion) *api.VersionedBlindedBeaconBlock {
	switch version {
	case spec.DataVersionBellatrix:
		return &api.VersionedBlindedBeaconBlock{
			Version:   version,
			Bellatrix: TestingBlindedBeaconBlock,
		}
	case spec.DataVersionCapella:
		return &api.VersionedBlindedBeaconBlock{
			Version: version,
			Capella: TestingBlindedBeaconBlockCapella,
		}
	default:
		panic("unsupported version")
	}
}

var TestingBlindedBeaconBlockBytesV = func(version spec.DataVersion) []byte {
	var ret []byte
	vBlk := TestingBlindedBeaconBlockV(version)

	switch version {
	case spec.DataVersionBellatrix:
		if vBlk.Bellatrix == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Bellatrix.MarshalSSZ()
	case spec.DataVersionCapella:
		if vBlk.Capella == nil {
			panic("empty block")
		}
		ret, _ = vBlk.Capella.MarshalSSZ()

	default:
		panic("unsupported version")
	}

	return ret
}

var TestingWrongBeaconBlockV = func(version spec.DataVersion) *spec.VersionedBeaconBlock {
	blkByts := TestingBeaconBlockBytesV(version)

	switch version {
	case spec.DataVersionBellatrix:
		ret := &bellatrix.BeaconBlock{}
		if err := ret.UnmarshalSSZ(blkByts); err != nil {
			panic(err.Error())
		}
		ret.Slot = 100
		return &spec.VersionedBeaconBlock{
			Version:   version,
			Bellatrix: ret,
		}
	case spec.DataVersionCapella:
		ret := &capella.BeaconBlock{}
		if err := ret.UnmarshalSSZ(blkByts); err != nil {
			panic(err.Error())
		}
		ret.Slot = TestingDutySlotCapella + 100
		return &spec.VersionedBeaconBlock{
			Version: version,
			Capella: ret,
		}

	default:
		panic("unsupported version")
	}
}

var TestingSignedBeaconBlockV = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	vBlk := TestingBeaconBlockV(version)

	switch version {
	case spec.DataVersionBellatrix:
		if vBlk.Bellatrix == nil {
			panic("empty block")
		}
		return &bellatrix.SignedBeaconBlock{
			Message:   vBlk.Bellatrix,
			Signature: signBeaconObject(vBlk.Bellatrix, types.DomainProposer, ks),
		}
	case spec.DataVersionCapella:
		if vBlk.Capella == nil {
			panic("empty block")
		}
		return &capella.SignedBeaconBlock{
			Message:   vBlk.Capella,
			Signature: signBeaconObject(vBlk.Capella, types.DomainProposer, ks),
		}

	default:
		panic("unsupported version")
	}
}

var TestingDutyEpochV = func(version spec.DataVersion) phase0.Epoch {
	switch version {
	case spec.DataVersionBellatrix:
		return TestingDutyEpochBellatrix
	case spec.DataVersionCapella:
		return TestingDutyEpochCapella

	default:
		panic("unsupported version")
	}
}

var TestingDutySlotV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionBellatrix:
		return TestingDutySlotBellatrix
	case spec.DataVersionCapella:
		return TestingDutySlotCapella

	default:
		panic("unsupported version")
	}
}

var VersionBySlot = func(slot phase0.Slot) spec.DataVersion {
	if slot < ForkEpochPraterCapella*32 {
		return spec.DataVersionBellatrix
	}
	return spec.DataVersionCapella
}

var TestingProposerDutyV = func(version spec.DataVersion) *types.Duty {
	duty := &types.Duty{
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

var TestingProposerDutyNextEpochV = func(version spec.DataVersion) *types.Duty {
	duty := &types.Duty{
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
	case spec.DataVersionBellatrix:
		duty.Slot = TestingDutySlotBellatrixNextEpoch
	case spec.DataVersionCapella:
		duty.Slot = TestingDutySlotCapellaNextEpoch

	default:
		panic("unsupported version")
	}

	return duty
}

var TestingInvalidDutySlotV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionBellatrix:
		return TestingDutySlotBellatrixInvalid
	case spec.DataVersionCapella:
		return TestingDutySlotCapellaInvalid

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
