package ssv

import (
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"
	apiv1fulu "github.com/attestantio/go-eth2-client/api/v1/fulu"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	denebspec "github.com/attestantio/go-eth2-client/spec/deneb"
	electraspec "github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
)

func TestEnsureBlindedProposalAlreadyBlinded(t *testing.T) {
	t.Run("capella", func(t *testing.T) {
		blindedBlock := &apiv1capella.BlindedBeaconBlock{}
		p := &api.VersionedProposal{Version: spec.DataVersionCapella, Blinded: true, CapellaBlinded: blindedBlock}

		blinded, marshaler, err := ensureBlindedProposal(p)
		require.NoError(t, err)
		require.Same(t, p, blinded)
		require.Same(t, blindedBlock, marshaler)
	})

	t.Run("deneb", func(t *testing.T) {
		blindedBlock := &apiv1deneb.BlindedBeaconBlock{}
		p := &api.VersionedProposal{Version: spec.DataVersionDeneb, Blinded: true, DenebBlinded: blindedBlock}

		blinded, marshaler, err := ensureBlindedProposal(p)
		require.NoError(t, err)
		require.Same(t, p, blinded)
		require.Same(t, blindedBlock, marshaler)
	})

	t.Run("electra", func(t *testing.T) {
		blindedBlock := &apiv1electra.BlindedBeaconBlock{}
		p := &api.VersionedProposal{Version: spec.DataVersionElectra, Blinded: true, ElectraBlinded: blindedBlock}

		blinded, marshaler, err := ensureBlindedProposal(p)
		require.NoError(t, err)
		require.Same(t, p, blinded)
		require.Same(t, blindedBlock, marshaler)
	})

	t.Run("fulu", func(t *testing.T) {
		blindedBlock := &apiv1electra.BlindedBeaconBlock{}
		p := &api.VersionedProposal{Version: spec.DataVersionFulu, Blinded: true, FuluBlinded: blindedBlock}

		blinded, marshaler, err := ensureBlindedProposal(p)
		require.NoError(t, err)
		require.Same(t, p, blinded)
		require.Same(t, blindedBlock, marshaler)
	})
}

func TestEnsureBlindedProposalFromFullBlock(t *testing.T) {
	t.Run("capella", func(t *testing.T) {
		full := capellaProposal()

		blinded, marshaler, err := ensureBlindedProposal(full)
		require.NoError(t, err)
		require.True(t, blinded.Blinded)
		require.Equal(t, spec.DataVersionCapella, blinded.Version)
		require.Same(t, blinded.CapellaBlinded, marshaler)

		txRoot, err := types.SSZTransactions(full.Capella.Body.ExecutionPayload.Transactions).HashTreeRoot()
		require.NoError(t, err)
		withdrawalsRoot, err := types.SSZWithdrawals(full.Capella.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		require.NoError(t, err)

		require.Equal(t, full.Capella.Slot, blinded.CapellaBlinded.Slot)
		require.Equal(t, full.Capella.ProposerIndex, blinded.CapellaBlinded.ProposerIndex)
		require.Equal(t, full.Capella.ParentRoot, blinded.CapellaBlinded.ParentRoot)
		require.Equal(t, full.Capella.StateRoot, blinded.CapellaBlinded.StateRoot)
		require.Equal(t, full.Capella.Body.Graffiti, blinded.CapellaBlinded.Body.Graffiti)
		require.Equal(t, phase0.Root(txRoot), blinded.CapellaBlinded.Body.ExecutionPayloadHeader.TransactionsRoot)
		require.Equal(t, phase0.Root(withdrawalsRoot), blinded.CapellaBlinded.Body.ExecutionPayloadHeader.WithdrawalsRoot)
	})

	t.Run("deneb", func(t *testing.T) {
		full := denebProposal()

		blinded, marshaler, err := ensureBlindedProposal(full)
		require.NoError(t, err)
		require.True(t, blinded.Blinded)
		require.Equal(t, spec.DataVersionDeneb, blinded.Version)
		require.Same(t, blinded.DenebBlinded, marshaler)

		txRoot, err := types.SSZTransactions(full.Deneb.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
		require.NoError(t, err)
		withdrawalsRoot, err := types.SSZWithdrawals(full.Deneb.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		require.NoError(t, err)

		require.Equal(t, full.Deneb.Block.Slot, blinded.DenebBlinded.Slot)
		require.Equal(t, full.Deneb.Block.ProposerIndex, blinded.DenebBlinded.ProposerIndex)
		require.Equal(t, full.Deneb.Block.ParentRoot, blinded.DenebBlinded.ParentRoot)
		require.Equal(t, full.Deneb.Block.StateRoot, blinded.DenebBlinded.StateRoot)
		require.Equal(t, full.Deneb.Block.Body.BlobKZGCommitments, blinded.DenebBlinded.Body.BlobKZGCommitments)
		require.Equal(t, phase0.Root(txRoot), blinded.DenebBlinded.Body.ExecutionPayloadHeader.TransactionsRoot)
		require.Equal(t, phase0.Root(withdrawalsRoot), blinded.DenebBlinded.Body.ExecutionPayloadHeader.WithdrawalsRoot)
		require.Equal(t, full.Deneb.Block.Body.ExecutionPayload.BlobGasUsed, blinded.DenebBlinded.Body.ExecutionPayloadHeader.BlobGasUsed)
		require.Equal(t, full.Deneb.Block.Body.ExecutionPayload.ExcessBlobGas, blinded.DenebBlinded.Body.ExecutionPayloadHeader.ExcessBlobGas)
	})

	t.Run("electra", func(t *testing.T) {
		full := electraProposal(spec.DataVersionElectra)

		blinded, marshaler, err := ensureBlindedProposal(full)
		require.NoError(t, err)
		require.True(t, blinded.Blinded)
		require.Equal(t, spec.DataVersionElectra, blinded.Version)
		require.Same(t, blinded.ElectraBlinded, marshaler)

		txRoot, err := types.SSZTransactions(full.Electra.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
		require.NoError(t, err)
		withdrawalsRoot, err := types.SSZWithdrawals(full.Electra.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		require.NoError(t, err)

		require.Equal(t, full.Electra.Block.Slot, blinded.ElectraBlinded.Slot)
		require.Equal(t, full.Electra.Block.ParentRoot, blinded.ElectraBlinded.ParentRoot)
		require.Equal(t, full.Electra.Block.Body.ExecutionRequests, blinded.ElectraBlinded.Body.ExecutionRequests)
		require.Equal(t, phase0.Root(txRoot), blinded.ElectraBlinded.Body.ExecutionPayloadHeader.TransactionsRoot)
		require.Equal(t, phase0.Root(withdrawalsRoot), blinded.ElectraBlinded.Body.ExecutionPayloadHeader.WithdrawalsRoot)
	})

	t.Run("fulu", func(t *testing.T) {
		full := electraProposal(spec.DataVersionFulu)

		blinded, marshaler, err := ensureBlindedProposal(full)
		require.NoError(t, err)
		require.True(t, blinded.Blinded)
		require.Equal(t, spec.DataVersionFulu, blinded.Version)
		require.Same(t, blinded.FuluBlinded, marshaler)

		txRoot, err := types.SSZTransactions(full.Fulu.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
		require.NoError(t, err)
		withdrawalsRoot, err := types.SSZWithdrawals(full.Fulu.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		require.NoError(t, err)

		require.Equal(t, full.Fulu.Block.Slot, blinded.FuluBlinded.Slot)
		require.Equal(t, full.Fulu.Block.ParentRoot, blinded.FuluBlinded.ParentRoot)
		require.Equal(t, full.Fulu.Block.Body.ExecutionRequests, blinded.FuluBlinded.Body.ExecutionRequests)
		require.Equal(t, phase0.Root(txRoot), blinded.FuluBlinded.Body.ExecutionPayloadHeader.TransactionsRoot)
		require.Equal(t, phase0.Root(withdrawalsRoot), blinded.FuluBlinded.Body.ExecutionPayloadHeader.WithdrawalsRoot)
	})
}

func capellaProposal() *api.VersionedProposal {
	payload := &capella.ExecutionPayload{
		ParentHash:    hash32(0x01),
		FeeRecipient:  executionAddress(0x02),
		StateRoot:     root32(0x03),
		ReceiptsRoot:  root32(0x04),
		LogsBloom:     bloom256(0x05),
		PrevRandao:    root32(0x06),
		BlockNumber:   7,
		GasLimit:      8,
		GasUsed:       9,
		Timestamp:     10,
		ExtraData:     []byte{0x0b, 0x0c},
		BaseFeePerGas: root32(0x0d),
		BlockHash:     hash32(0x0e),
		Transactions:  []bellatrix.Transaction{{0x0f, 0x10}},
		Withdrawals: []*capella.Withdrawal{{
			Index:          11,
			ValidatorIndex: 12,
			Address:        executionAddress(0x11),
			Amount:         13,
		}},
	}

	return &api.VersionedProposal{
		Version: spec.DataVersionCapella,
		Capella: &capella.BeaconBlock{
			Slot:          1,
			ProposerIndex: 2,
			ParentRoot:    root32(0x12),
			StateRoot:     root32(0x13),
			Body: &capella.BeaconBlockBody{
				ETH1Data:         eth1Data(0x14),
				Graffiti:         root32(0x15),
				SyncAggregate:    &altair.SyncAggregate{},
				ExecutionPayload: payload,
			},
		},
	}
}

func denebProposal() *api.VersionedProposal {
	payload := &denebspec.ExecutionPayload{
		ParentHash:    hash32(0x21),
		FeeRecipient:  executionAddress(0x22),
		StateRoot:     root32(0x23),
		ReceiptsRoot:  root32(0x24),
		LogsBloom:     bloom256(0x25),
		PrevRandao:    root32(0x26),
		BlockNumber:   27,
		GasLimit:      28,
		GasUsed:       29,
		Timestamp:     30,
		ExtraData:     []byte{0x2b, 0x2c},
		BaseFeePerGas: uint256.NewInt(31),
		BlockHash:     hash32(0x2d),
		Transactions:  []bellatrix.Transaction{{0x2e, 0x2f}},
		Withdrawals: []*capella.Withdrawal{{
			Index:          32,
			ValidatorIndex: 33,
			Address:        executionAddress(0x30),
			Amount:         34,
		}},
		BlobGasUsed:   35,
		ExcessBlobGas: 36,
	}

	contents := &apiv1deneb.BlockContents{
		Block: &denebspec.BeaconBlock{
			Slot:          37,
			ProposerIndex: 38,
			ParentRoot:    root32(0x31),
			StateRoot:     root32(0x32),
			Body: &denebspec.BeaconBlockBody{
				ETH1Data:           eth1Data(0x33),
				Graffiti:           root32(0x34),
				SyncAggregate:      &altair.SyncAggregate{},
				ExecutionPayload:   payload,
				BlobKZGCommitments: []denebspec.KZGCommitment{kzgCommitment(0x35)},
			},
		},
		KZGProofs: []denebspec.KZGProof{kzgProof(0x36)},
		Blobs:     []denebspec.Blob{blob(0x37)},
	}

	return &api.VersionedProposal{
		Version: spec.DataVersionDeneb,
		Deneb:   contents,
	}
}

func electraProposal(version spec.DataVersion) *api.VersionedProposal {
	payload := &denebspec.ExecutionPayload{
		ParentHash:    hash32(0x41),
		FeeRecipient:  executionAddress(0x42),
		StateRoot:     root32(0x43),
		ReceiptsRoot:  root32(0x44),
		LogsBloom:     bloom256(0x45),
		PrevRandao:    root32(0x46),
		BlockNumber:   47,
		GasLimit:      48,
		GasUsed:       49,
		Timestamp:     50,
		ExtraData:     []byte{0x4b, 0x4c},
		BaseFeePerGas: uint256.NewInt(51),
		BlockHash:     hash32(0x4d),
		Transactions:  []bellatrix.Transaction{{0x4e, 0x4f}},
		Withdrawals: []*capella.Withdrawal{{
			Index:          52,
			ValidatorIndex: 53,
			Address:        executionAddress(0x50),
			Amount:         54,
		}},
		BlobGasUsed:   55,
		ExcessBlobGas: 56,
	}

	contents := &apiv1electra.BlockContents{
		Block: &electraspec.BeaconBlock{
			Slot:          57,
			ProposerIndex: 58,
			ParentRoot:    root32(0x51),
			StateRoot:     root32(0x52),
			Body: &electraspec.BeaconBlockBody{
				ETH1Data:           eth1Data(0x53),
				Graffiti:           root32(0x54),
				SyncAggregate:      &altair.SyncAggregate{},
				ExecutionPayload:   payload,
				BlobKZGCommitments: []denebspec.KZGCommitment{kzgCommitment(0x55)},
				ExecutionRequests:  &electraspec.ExecutionRequests{},
			},
		},
		KZGProofs: []denebspec.KZGProof{kzgProof(0x56)},
		Blobs:     []denebspec.Blob{blob(0x57)},
	}

	if version == spec.DataVersionFulu {
		return &api.VersionedProposal{
			Version: version,
			Fulu: &apiv1fulu.BlockContents{
				Block:     contents.Block,
				KZGProofs: contents.KZGProofs,
				Blobs:     contents.Blobs,
			},
		}
	}

	return &api.VersionedProposal{
		Version: version,
		Electra: contents,
	}
}

func eth1Data(fill byte) *phase0.ETH1Data {
	return &phase0.ETH1Data{
		DepositRoot:  root32(fill),
		DepositCount: uint64(fill),
		BlockHash:    filledBytes(fill, 32),
	}
}

func root32(fill byte) [32]byte {
	var out [32]byte
	for i := range out {
		out[i] = fill
	}
	return out
}

func hash32(fill byte) phase0.Hash32 {
	return phase0.Hash32(root32(fill))
}

func executionAddress(fill byte) bellatrix.ExecutionAddress {
	var out bellatrix.ExecutionAddress
	for i := range out {
		out[i] = fill
	}
	return out
}

func bloom256(fill byte) [256]byte {
	var out [256]byte
	for i := range out {
		out[i] = fill
	}
	return out
}

func kzgCommitment(fill byte) denebspec.KZGCommitment {
	var out denebspec.KZGCommitment
	for i := range out {
		out[i] = fill
	}
	return out
}

func kzgProof(fill byte) denebspec.KZGProof {
	var out denebspec.KZGProof
	for i := range out {
		out[i] = fill
	}
	return out
}

func blob(fill byte) denebspec.Blob {
	var out denebspec.Blob
	for i := range out {
		out[i] = fill
	}
	return out
}

func filledBytes(fill byte, count int) []byte {
	out := make([]byte, count)
	for i := range out {
		out[i] = fill
	}
	return out
}
