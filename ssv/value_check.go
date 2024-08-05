package ssv

import (
	"bytes"
	"math"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// BeaconVoteValueCheckF creates a value check function for a beacon vote
// It filters out the slashable public keys from the shares and creates a value check function
// that checks if the beacon vote is slashable by the filtered public keys
// It also does sanity checks on the beacon vote data
// It returns an error if shares assigned to the duty are missing
func BeaconVoteValueCheckF(duty *types.CommitteeDuty, signer types.BeaconSigner, beaconNetwork types.BeaconNetwork,
	availableShares map[phase0.ValidatorIndex]*types.Share) (qbft.ProposedValueCheckF, error) {
	duties := duty.ValidatorDuties
	slashablePublicKeys := []types.ShareValidatorPK{}
	for _, duty := range duties {
		if duty.Type == types.BNRoleAttester {
			share, ok := availableShares[duty.ValidatorIndex]
			if !ok {
				return nil, errors.New("assigned validator duty doesn't have a validator share")
			}
			slashablePublicKeys = append(slashablePublicKeys, share.SharePubKey)
		}
	}
	valueCheckF := createBeaconVoteValueCheckF(signer, duty.DutySlot(), slashablePublicKeys,
		beaconNetwork.EstimatedCurrentEpoch())
	return valueCheckF, nil
}

func dutyValueCheck(
	duty *types.ValidatorDuty,
	network types.BeaconNetwork,
	expectedType types.BeaconRole,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
) error {
	if network.EstimatedEpochAtSlot(duty.Slot) > network.EstimatedCurrentEpoch()+1 {
		return errors.New("duty epoch is into far future")
	}

	if expectedType != duty.Type {
		return errors.New("wrong beacon role type")
	}

	if !bytes.Equal(validatorPK[:], duty.PubKey[:]) {
		return errors.New("wrong validator pk")
	}

	if validatorIndex != duty.ValidatorIndex {
		return errors.New("wrong validator index")
	}

	return nil
}

func createBeaconVoteValueCheckF(
	signer types.BeaconSigner,
	slot phase0.Slot,
	sharePublicKeys []types.ShareValidatorPK,
	estimatedCurrentEpoch phase0.Epoch,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		bv := types.BeaconVote{}
		if err := bv.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding beacon vote")
		}

		if bv.Target.Epoch > estimatedCurrentEpoch+1 {
			return errors.New("attestation data target epoch is into far future")
		}

		if bv.Source.Epoch >= bv.Target.Epoch {
			return errors.New("attestation data source >= target")
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

		if blockData, _, err := cd.GetBlindedBlockData(); err == nil {
			slot, err := blockData.Slot()
			if err != nil {
				return errors.Wrap(err, "failed to get slot from blinded block data")
			}
			return signer.IsBeaconBlockSlashable(sharePublicKey, slot)
		}
		if blockData, _, err := cd.GetBlockData(); err == nil {
			slot, err := blockData.Slot()
			if err != nil {
				return errors.Wrap(err, "failed to get slot from block data")
			}
			return signer.IsBeaconBlockSlashable(sharePublicKey, slot)
		}

		return errors.New("no block data")
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
