package ssv

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types"
)

// ensureBlindedProposal returns a blinded proposal and the concrete SSZ-marshaler
// that must be encoded into proposer consensus data. If the input is already
// blinded, it is returned unchanged.
func ensureBlindedProposal(p *api.VersionedProposal) (*api.VersionedProposal, ssz.Marshaler, error) {
	if p == nil {
		return nil, nil, fmt.Errorf("nil proposal")
	}

	if p.Blinded {
		switch p.Version {
		case spec.DataVersionCapella:
			if p.CapellaBlinded == nil {
				return nil, nil, fmt.Errorf("capella blinded block is nil")
			}
			return p, p.CapellaBlinded, nil
		case spec.DataVersionDeneb:
			if p.DenebBlinded == nil {
				return nil, nil, fmt.Errorf("deneb blinded block is nil")
			}
			return p, p.DenebBlinded, nil
		case spec.DataVersionElectra:
			if p.ElectraBlinded == nil {
				return nil, nil, fmt.Errorf("electra blinded block is nil")
			}
			return p, p.ElectraBlinded, nil
		case spec.DataVersionFulu:
			if p.FuluBlinded == nil {
				return nil, nil, fmt.Errorf("fulu blinded block is nil")
			}
			return p, p.FuluBlinded, nil
		default:
			return nil, nil, fmt.Errorf("unsupported proposal version %d", p.Version)
		}
	}

	switch p.Version {
	case spec.DataVersionCapella:
		if p.Capella == nil || p.Capella.Body == nil || p.Capella.Body.ExecutionPayload == nil {
			return nil, nil, fmt.Errorf("capella block or payload is nil")
		}

		txRoot, err := types.SSZTransactions(p.Capella.Body.ExecutionPayload.Transactions).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute capella transactions root: %w", err)
		}
		withdrawalsRoot, err := types.SSZWithdrawals(p.Capella.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute capella withdrawals root: %w", err)
		}

		blinded := &apiv1capella.BlindedBeaconBlock{
			Slot:          p.Capella.Slot,
			ProposerIndex: p.Capella.ProposerIndex,
			ParentRoot:    p.Capella.ParentRoot,
			StateRoot:     p.Capella.StateRoot,
			Body: &apiv1capella.BlindedBeaconBlockBody{
				RANDAOReveal:      p.Capella.Body.RANDAOReveal,
				ETH1Data:          p.Capella.Body.ETH1Data,
				Graffiti:          p.Capella.Body.Graffiti,
				ProposerSlashings: p.Capella.Body.ProposerSlashings,
				AttesterSlashings: p.Capella.Body.AttesterSlashings,
				Attestations:      p.Capella.Body.Attestations,
				Deposits:          p.Capella.Body.Deposits,
				VoluntaryExits:    p.Capella.Body.VoluntaryExits,
				SyncAggregate:     p.Capella.Body.SyncAggregate,
				ExecutionPayloadHeader: &capella.ExecutionPayloadHeader{
					ParentHash:       p.Capella.Body.ExecutionPayload.ParentHash,
					FeeRecipient:     p.Capella.Body.ExecutionPayload.FeeRecipient,
					StateRoot:        p.Capella.Body.ExecutionPayload.StateRoot,
					ReceiptsRoot:     p.Capella.Body.ExecutionPayload.ReceiptsRoot,
					LogsBloom:        p.Capella.Body.ExecutionPayload.LogsBloom,
					PrevRandao:       p.Capella.Body.ExecutionPayload.PrevRandao,
					BlockNumber:      p.Capella.Body.ExecutionPayload.BlockNumber,
					GasLimit:         p.Capella.Body.ExecutionPayload.GasLimit,
					GasUsed:          p.Capella.Body.ExecutionPayload.GasUsed,
					Timestamp:        p.Capella.Body.ExecutionPayload.Timestamp,
					ExtraData:        p.Capella.Body.ExecutionPayload.ExtraData,
					BaseFeePerGas:    p.Capella.Body.ExecutionPayload.BaseFeePerGas,
					BlockHash:        p.Capella.Body.ExecutionPayload.BlockHash,
					TransactionsRoot: txRoot,
					WithdrawalsRoot:  withdrawalsRoot,
				},
				BLSToExecutionChanges: p.Capella.Body.BLSToExecutionChanges,
			},
		}

		return &api.VersionedProposal{Version: p.Version, Blinded: true, CapellaBlinded: blinded}, blinded, nil
	case spec.DataVersionDeneb:
		if p.Deneb == nil || p.Deneb.Block == nil || p.Deneb.Block.Body == nil || p.Deneb.Block.Body.ExecutionPayload == nil {
			return nil, nil, fmt.Errorf("deneb block or payload is nil")
		}

		txRoot, err := types.SSZTransactions(p.Deneb.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute deneb transactions root: %w", err)
		}
		withdrawalsRoot, err := types.SSZWithdrawals(p.Deneb.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute deneb withdrawals root: %w", err)
		}

		blinded := &apiv1deneb.BlindedBeaconBlock{
			Slot:          p.Deneb.Block.Slot,
			ProposerIndex: p.Deneb.Block.ProposerIndex,
			ParentRoot:    p.Deneb.Block.ParentRoot,
			StateRoot:     p.Deneb.Block.StateRoot,
			Body: &apiv1deneb.BlindedBeaconBlockBody{
				RANDAOReveal:      p.Deneb.Block.Body.RANDAOReveal,
				ETH1Data:          p.Deneb.Block.Body.ETH1Data,
				Graffiti:          p.Deneb.Block.Body.Graffiti,
				ProposerSlashings: p.Deneb.Block.Body.ProposerSlashings,
				AttesterSlashings: p.Deneb.Block.Body.AttesterSlashings,
				Attestations:      p.Deneb.Block.Body.Attestations,
				Deposits:          p.Deneb.Block.Body.Deposits,
				VoluntaryExits:    p.Deneb.Block.Body.VoluntaryExits,
				SyncAggregate:     p.Deneb.Block.Body.SyncAggregate,
				ExecutionPayloadHeader: &deneb.ExecutionPayloadHeader{
					ParentHash:       p.Deneb.Block.Body.ExecutionPayload.ParentHash,
					FeeRecipient:     p.Deneb.Block.Body.ExecutionPayload.FeeRecipient,
					StateRoot:        p.Deneb.Block.Body.ExecutionPayload.StateRoot,
					ReceiptsRoot:     p.Deneb.Block.Body.ExecutionPayload.ReceiptsRoot,
					LogsBloom:        p.Deneb.Block.Body.ExecutionPayload.LogsBloom,
					PrevRandao:       p.Deneb.Block.Body.ExecutionPayload.PrevRandao,
					BlockNumber:      p.Deneb.Block.Body.ExecutionPayload.BlockNumber,
					GasLimit:         p.Deneb.Block.Body.ExecutionPayload.GasLimit,
					GasUsed:          p.Deneb.Block.Body.ExecutionPayload.GasUsed,
					Timestamp:        p.Deneb.Block.Body.ExecutionPayload.Timestamp,
					ExtraData:        p.Deneb.Block.Body.ExecutionPayload.ExtraData,
					BaseFeePerGas:    p.Deneb.Block.Body.ExecutionPayload.BaseFeePerGas,
					BlockHash:        p.Deneb.Block.Body.ExecutionPayload.BlockHash,
					TransactionsRoot: txRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      p.Deneb.Block.Body.ExecutionPayload.BlobGasUsed,
					ExcessBlobGas:    p.Deneb.Block.Body.ExecutionPayload.ExcessBlobGas,
				},
				BLSToExecutionChanges: p.Deneb.Block.Body.BLSToExecutionChanges,
				BlobKZGCommitments:    p.Deneb.Block.Body.BlobKZGCommitments,
			},
		}

		return &api.VersionedProposal{Version: p.Version, Blinded: true, DenebBlinded: blinded}, blinded, nil
	case spec.DataVersionElectra:
		if p.Electra == nil || p.Electra.Block == nil || p.Electra.Block.Body == nil || p.Electra.Block.Body.ExecutionPayload == nil {
			return nil, nil, fmt.Errorf("electra block or payload is nil")
		}

		txRoot, err := types.SSZTransactions(p.Electra.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute electra transactions root: %w", err)
		}
		withdrawalsRoot, err := types.SSZWithdrawals(p.Electra.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute electra withdrawals root: %w", err)
		}

		blinded := &apiv1electra.BlindedBeaconBlock{
			Slot:          p.Electra.Block.Slot,
			ProposerIndex: p.Electra.Block.ProposerIndex,
			ParentRoot:    p.Electra.Block.ParentRoot,
			StateRoot:     p.Electra.Block.StateRoot,
			Body: &apiv1electra.BlindedBeaconBlockBody{
				RANDAOReveal:      p.Electra.Block.Body.RANDAOReveal,
				ETH1Data:          p.Electra.Block.Body.ETH1Data,
				Graffiti:          p.Electra.Block.Body.Graffiti,
				ProposerSlashings: p.Electra.Block.Body.ProposerSlashings,
				AttesterSlashings: p.Electra.Block.Body.AttesterSlashings,
				Attestations:      p.Electra.Block.Body.Attestations,
				Deposits:          p.Electra.Block.Body.Deposits,
				VoluntaryExits:    p.Electra.Block.Body.VoluntaryExits,
				SyncAggregate:     p.Electra.Block.Body.SyncAggregate,
				ExecutionPayloadHeader: &deneb.ExecutionPayloadHeader{
					ParentHash:       p.Electra.Block.Body.ExecutionPayload.ParentHash,
					FeeRecipient:     p.Electra.Block.Body.ExecutionPayload.FeeRecipient,
					StateRoot:        p.Electra.Block.Body.ExecutionPayload.StateRoot,
					ReceiptsRoot:     p.Electra.Block.Body.ExecutionPayload.ReceiptsRoot,
					LogsBloom:        p.Electra.Block.Body.ExecutionPayload.LogsBloom,
					PrevRandao:       p.Electra.Block.Body.ExecutionPayload.PrevRandao,
					BlockNumber:      p.Electra.Block.Body.ExecutionPayload.BlockNumber,
					GasLimit:         p.Electra.Block.Body.ExecutionPayload.GasLimit,
					GasUsed:          p.Electra.Block.Body.ExecutionPayload.GasUsed,
					Timestamp:        p.Electra.Block.Body.ExecutionPayload.Timestamp,
					ExtraData:        p.Electra.Block.Body.ExecutionPayload.ExtraData,
					BaseFeePerGas:    p.Electra.Block.Body.ExecutionPayload.BaseFeePerGas,
					BlockHash:        p.Electra.Block.Body.ExecutionPayload.BlockHash,
					TransactionsRoot: txRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      p.Electra.Block.Body.ExecutionPayload.BlobGasUsed,
					ExcessBlobGas:    p.Electra.Block.Body.ExecutionPayload.ExcessBlobGas,
				},
				BLSToExecutionChanges: p.Electra.Block.Body.BLSToExecutionChanges,
				BlobKZGCommitments:    p.Electra.Block.Body.BlobKZGCommitments,
				ExecutionRequests:     p.Electra.Block.Body.ExecutionRequests,
			},
		}

		return &api.VersionedProposal{Version: p.Version, Blinded: true, ElectraBlinded: blinded}, blinded, nil
	case spec.DataVersionFulu:
		if p.Fulu == nil || p.Fulu.Block == nil || p.Fulu.Block.Body == nil || p.Fulu.Block.Body.ExecutionPayload == nil {
			return nil, nil, fmt.Errorf("fulu block or payload is nil")
		}

		txRoot, err := types.SSZTransactions(p.Fulu.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute fulu transactions root: %w", err)
		}
		withdrawalsRoot, err := types.SSZWithdrawals(p.Fulu.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute fulu withdrawals root: %w", err)
		}

		blinded := &apiv1electra.BlindedBeaconBlock{
			Slot:          p.Fulu.Block.Slot,
			ProposerIndex: p.Fulu.Block.ProposerIndex,
			ParentRoot:    p.Fulu.Block.ParentRoot,
			StateRoot:     p.Fulu.Block.StateRoot,
			Body: &apiv1electra.BlindedBeaconBlockBody{
				RANDAOReveal:      p.Fulu.Block.Body.RANDAOReveal,
				ETH1Data:          p.Fulu.Block.Body.ETH1Data,
				Graffiti:          p.Fulu.Block.Body.Graffiti,
				ProposerSlashings: p.Fulu.Block.Body.ProposerSlashings,
				AttesterSlashings: p.Fulu.Block.Body.AttesterSlashings,
				Attestations:      p.Fulu.Block.Body.Attestations,
				Deposits:          p.Fulu.Block.Body.Deposits,
				VoluntaryExits:    p.Fulu.Block.Body.VoluntaryExits,
				SyncAggregate:     p.Fulu.Block.Body.SyncAggregate,
				ExecutionPayloadHeader: &deneb.ExecutionPayloadHeader{
					ParentHash:       p.Fulu.Block.Body.ExecutionPayload.ParentHash,
					FeeRecipient:     p.Fulu.Block.Body.ExecutionPayload.FeeRecipient,
					StateRoot:        p.Fulu.Block.Body.ExecutionPayload.StateRoot,
					ReceiptsRoot:     p.Fulu.Block.Body.ExecutionPayload.ReceiptsRoot,
					LogsBloom:        p.Fulu.Block.Body.ExecutionPayload.LogsBloom,
					PrevRandao:       p.Fulu.Block.Body.ExecutionPayload.PrevRandao,
					BlockNumber:      p.Fulu.Block.Body.ExecutionPayload.BlockNumber,
					GasLimit:         p.Fulu.Block.Body.ExecutionPayload.GasLimit,
					GasUsed:          p.Fulu.Block.Body.ExecutionPayload.GasUsed,
					Timestamp:        p.Fulu.Block.Body.ExecutionPayload.Timestamp,
					ExtraData:        p.Fulu.Block.Body.ExecutionPayload.ExtraData,
					BaseFeePerGas:    p.Fulu.Block.Body.ExecutionPayload.BaseFeePerGas,
					BlockHash:        p.Fulu.Block.Body.ExecutionPayload.BlockHash,
					TransactionsRoot: txRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      p.Fulu.Block.Body.ExecutionPayload.BlobGasUsed,
					ExcessBlobGas:    p.Fulu.Block.Body.ExecutionPayload.ExcessBlobGas,
				},
				BLSToExecutionChanges: p.Fulu.Block.Body.BLSToExecutionChanges,
				BlobKZGCommitments:    p.Fulu.Block.Body.BlobKZGCommitments,
				ExecutionRequests:     p.Fulu.Block.Body.ExecutionRequests,
			},
		}

		return &api.VersionedProposal{Version: p.Version, Blinded: true, FuluBlinded: blinded}, blinded, nil
	default:
		return nil, nil, fmt.Errorf("unsupported proposal version %d", p.Version)
	}
}
