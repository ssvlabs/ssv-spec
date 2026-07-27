package ssv

import (
	"bytes"
	"fmt"
	"math"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

func dutyValueCheck(
	duty *types.ValidatorDuty,
	network types.BeaconNetwork,
	expectedType types.BeaconRole,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
) error {
	if network.EstimatedEpochAtSlot(duty.Slot) > network.EstimatedCurrentEpoch()+1 {
		return types.NewError(types.DutyEpochTooFarFutureErrorCode, "duty epoch is into far future")
	}

	if expectedType != duty.Type {
		return types.NewError(types.WrongBeaconRoleTypeErrorCode, "wrong beacon role type")
	}

	if !bytes.Equal(validatorPK[:], duty.PubKey[:]) {
		return types.NewError(types.WrongValidatorPubkeyErrorCode, "wrong validator pk")
	}

	if validatorIndex != duty.ValidatorIndex {
		return types.NewError(types.WrongValidatorIndexErrorCode, "wrong validator index")
	}

	return nil
}

func BeaconVoteValueCheckF(
	signer types.BeaconSigner,
	slot phase0.Slot,
	sharePublicKeys []types.ShareValidatorPK,
	expectedSource phase0.Epoch,
	expectedTarget phase0.Epoch,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		bv := types.BeaconVote{}
		if err := bv.Decode(data); err != nil {
			return types.WrapError(types.DecodeBeaconVoteErrorCode, fmt.Errorf("failed decoding beacon vote: %w", err))
		}

		if err := bv.Validate(); err != nil {
			return err
		}

		if bv.Source.Epoch != expectedSource {
			return types.NewError(types.CheckpointMismatch,
				fmt.Sprintf("attestation data source checkpoint %d does not match expected %d",
					bv.Source.Epoch, expectedSource))
		}

		if bv.Target.Epoch != expectedTarget {
			return types.NewError(types.CheckpointMismatch,
				fmt.Sprintf("attestation data target checkpoint %d does not match expected %d",
					bv.Target.Epoch, expectedTarget))
		}

		attestationData := &phase0.AttestationData{
			Slot: slot,
			// Consensus data is unaware of CommitteeIndex
			// We use -1 to not run into issues with the duplicate value slashing check:
			// (data_1 != data_2 and data_1.target.epoch == data_2.target.epoch)
			Index:           math.MaxUint64,
			BeaconBlockRoot: bv.BlockRoot,
			Source:          bv.Source,
			Target:          bv.Target,
		}

		for _, sharePublicKey := range sharePublicKeys {
			if err := signer.IsAttestationSlashable(sharePublicKey, attestationData); err != nil {
				return err
			}
		}
		return nil
	}
}

// GloasBeaconVoteValueCheckF is the Gloas (ePBS) variant of BeaconVoteValueCheckF (SIP #94 §2). It
// mirrors the checkpoint/slashability checks and adds two Gloas rules: reject AttestationDataIndex > 1
// (only 0 = payload-absent, 1 = payload-present are valid), and build the slashability AttestationData
// with the decided index rather than the math.MaxUint64 sentinel, so cross-index equivocation over the
// same (source, target, slot) trips IsAttestationSlashable.
func GloasBeaconVoteValueCheckF(
	signer types.BeaconSigner,
	slot phase0.Slot,
	sharePublicKeys []types.ShareValidatorPK,
	expectedSource phase0.Epoch,
	expectedTarget phase0.Epoch,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		bv := types.GloasBeaconVote{}
		if err := bv.Decode(data); err != nil {
			return types.WrapError(types.DecodeGloasBeaconVoteErrorCode, fmt.Errorf("failed decoding gloas beacon vote: %w", err))
		}

		if bv.AttestationDataIndex > 1 {
			return types.NewError(types.GloasBeaconVoteInvalidIndexErrorCode,
				fmt.Sprintf("attestation data index %d must be 0 or 1", bv.AttestationDataIndex))
		}

		if bv.Source.Epoch >= bv.Target.Epoch {
			return types.NewError(types.AttestationSourceNotLessThanTargetErrorCode, "attestation data source >= target")
		}

		if bv.Source.Epoch != expectedSource {
			return types.NewError(types.CheckpointMismatch,
				fmt.Sprintf("attestation data source checkpoint %d does not match expected %d",
					bv.Source.Epoch, expectedSource))
		}

		if bv.Target.Epoch != expectedTarget {
			return types.NewError(types.CheckpointMismatch,
				fmt.Sprintf("attestation data target checkpoint %d does not match expected %d",
					bv.Target.Epoch, expectedTarget))
		}

		attestationData := &phase0.AttestationData{
			Slot: slot,
			// The 0/1 index is a meaningful part of the Gloas vote, so it goes into the slashability data
			// directly (not the pre-Gloas math.MaxUint64 sentinel) — cross-index equivocation over the
			// same (source, target, slot) is then double-vote slashable.
			Index:           bv.AttestationDataIndex,
			BeaconBlockRoot: bv.BlockRoot,
			Source:          bv.Source,
			Target:          bv.Target,
		}

		for _, sharePublicKey := range sharePublicKeys {
			if err := signer.IsAttestationSlashable(sharePublicKey, attestationData); err != nil {
				return err
			}
		}
		return nil
	}
}

func ProposerValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
	sharePublicKey []byte,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.ProposerConsensusData{}
		if err := cd.Decode(data); err != nil {
			return types.WrapError(types.ProposerConsensusDataDecodeErrorCode, errors.Wrap(err, "failed decoding consensus data"))
		}
		// Gloas (ePBS §4): the block is opaque to the types layer — GetBlockData()/Validate() have no
		// Gloas arm (go-eth2-client's api.VersionedProposal can't carry Gloas). Branch on the fork
		// before any type-layer validation and decode the block here, so a Gloas value never routes
		// through Validate()/GetBlockData(). Pre-Gloas is unchanged.
		// Exact match, not >=: a future fork's consensus data must not be decoded as a Gloas block —
		// it falls through to Validate() and errors as an unknown version.
		if cd.Version == gloas.DataVersionGloas {
			if err := dutyValueCheck(&cd.Duty, network, types.BNRoleProposer, validatorPK, validatorIndex); err != nil {
				return errors.Wrap(err, "duty invalid")
			}
			block, err := gloas.DecodeBeaconBlock(cd.DataSSZ)
			if err != nil {
				return types.WrapError(types.UnmarshalSSZErrorCode, errors.Wrap(err, "failed decoding gloas beacon block"))
			}
			// The QBFT-agreed block must be for the duty's slot; without this pin the cluster could
			// agree on a block for a different slot (SIP #94 §4). Mirrors §6's duty-slot match.
			if block.Slot != cd.Duty.Slot {
				return types.NewError(types.ProposerBlockSlotMismatchErrorCode, "gloas block slot does not match duty slot")
			}
			return signer.IsBeaconBlockSlashable(sharePublicKey, block.Slot)
		}

		if err := cd.Validate(); err != nil {
			return types.NewError(types.QBFTValueInvalidErrorCode, fmt.Sprintf("invalid value: %v", err.Error()))
		}

		if err := dutyValueCheck(&cd.Duty, network, types.BNRoleProposer, validatorPK, validatorIndex); err != nil {
			return errors.Wrap(err, "duty invalid")
		}

		blockData, _, err := cd.GetBlockData()
		if err != nil {
			return errors.Wrap(err, "could not get block data")
		}
		slot, err := blockData.Slot()
		if err != nil {
			return errors.Wrap(err, "failed to get slot from block data")
		}
		return signer.IsBeaconBlockSlashable(sharePublicKey, slot)
	}
}

func AggregatorCommitteeValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.AggregatorCommitteeConsensusData{}
		if err := cd.Decode(data); err != nil {
			return types.WrapError(types.AggCommConsensusDataDecodeErrorCode, errors.Wrap(err, "failed decoding aggregator committee consensus data"))
		}
		if err := cd.Validate(); err != nil {
			return errors.Wrap(err, "invalid value")
		}

		return nil
	}
}
