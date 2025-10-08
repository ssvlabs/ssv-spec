package ssv

import (
	"bytes"
	"fmt"
	"math"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
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
	estimatedCurrentEpoch phase0.Epoch,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		bv := types.BeaconVote{}
		if err := bv.Decode(data); err != nil {
			return types.NewError(types.DecodeBeaconVoteErrorCode, fmt.Sprintf("failed decoding beacon vote: %v", err))
		}

		if bv.Target.Epoch > estimatedCurrentEpoch+1 {
			return types.NewError(types.AttestationTargetEpochTooFarFutureErrorCode, "attestation data target epoch is into far future")
		}

		if bv.Source.Epoch >= bv.Target.Epoch {
			return types.NewError(types.AttestationSourceNotLessThanTargetErrorCode, "attestation data source >= target")
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

func ProposerValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
	sharePublicKey []byte,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.ValidatorConsensusData{}
		if err := cd.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding consensus data")
		}
		if err := cd.Validate(); err != nil {
			return errors.Wrap(err, "invalid value")
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

func AggregatorValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.ValidatorConsensusData{}
		if err := cd.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding consensus data")
		}
		if err := cd.Validate(); err != nil {
			return errors.Wrap(err, "invalid value")
		}

		if err := dutyValueCheck(&cd.Duty, network, types.BNRoleAggregator, validatorPK, validatorIndex); err != nil {
			return errors.Wrap(err, "duty invalid")
		}
		return nil
	}
}

func SyncCommitteeContributionValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.ValidatorConsensusData{}
		if err := cd.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding consensus data")
		}
		if err := cd.Validate(); err != nil {
			return errors.Wrap(err, "invalid value")
		}

		if err := dutyValueCheck(&cd.Duty, network, types.BNRoleSyncCommitteeContribution, validatorPK, validatorIndex); err != nil {
			return errors.Wrap(err, "duty invalid")
		}

		//contributions, _ := cd.GetSyncCommitteeContributions()
		//
		//for _, c := range contributions {
		//	// TODO check we have selection proof for contribution
		//	// TODO check slot == duty slot
		//	// TODO check beacon block root somehow? maybe all beacon block roots should be equal?
		//
		//}
		return nil
	}
}
