package ssv

import (
	"bytes"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

func dutyValueCheck(
	duty *types.Duty,
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

	if !bytes.Equal(validatorPK, duty.PubKey[:]) {
		return errors.New("wrong validator pk")
	}

	if validatorIndex != duty.ValidatorIndex {
		return errors.New("wrong validator index")
	}

	return nil
}

func AttesterValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.ConsensusData{}
		if err := cd.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding consensus data")
		}

		if err := dutyValueCheck(cd.Duty, network, types.BNRoleAttester, validatorPK, validatorIndex); err != nil {
			return errors.Wrap(err, "duty invalid")
		}

		if cd.AttestationData == nil {
			return errors.New("attestation data nil")
		}

		if cd.Duty.Slot != cd.AttestationData.Slot {
			return errors.New("attestation data slot != duty slot")
		}

		if cd.Duty.CommitteeIndex != cd.AttestationData.Index {
			return errors.New("attestation data CommitteeIndex != duty CommitteeIndex")
		}

		if cd.AttestationData.Target.Epoch > network.EstimatedCurrentEpoch()+1 {
			return errors.New("attestation data target epoch is into far future")
		}

		if cd.AttestationData.Source.Epoch >= cd.AttestationData.Target.Epoch {
			return errors.New("attestation data source > target")
		}

		return signer.IsAttestationSlashable(cd.AttestationData)
	}
}

func ProposerValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.ConsensusData{}
		if err := cd.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding consensus data")
		}

		if err := dutyValueCheck(cd.Duty, network, types.BNRoleProposer, validatorPK, validatorIndex); err != nil {
			return errors.Wrap(err, "duty invalid")
		}
		return nil
	}
}

func AggregatorValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.ConsensusData{}
		if err := cd.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding consensus data")
		}

		if err := dutyValueCheck(cd.Duty, network, types.BNRoleAggregator, validatorPK, validatorIndex); err != nil {
			return errors.Wrap(err, "duty invalid")
		}
		return nil
	}
}

func SyncCommitteeValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.ConsensusData{}
		if err := cd.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding consensus data")
		}

		if err := dutyValueCheck(cd.Duty, network, types.BNRoleSyncCommittee, validatorPK, validatorIndex); err != nil {
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
		cd := &types.ConsensusData{}
		if err := cd.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding consensus data")
		}

		if err := dutyValueCheck(cd.Duty, network, types.BNRoleSyncCommitteeContribution, validatorPK, validatorIndex); err != nil {
			return errors.Wrap(err, "duty invalid")
		}

		for _, c := range cd.SyncCommitteeContribution {
			// nolint
			if c.Slot == 0 {
				// TODO - can remove
			}

			// TODO check we have selection proof for contribution
			// TODO check slot == duty slot
			// TODO check beacon block root somehow? maybe all beacon block roots should be equal?

		}
		return nil
	}
}
