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

func BeaconAttestationValueCheck(
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

func BeaconBlockValueCheck(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		return nil
	}
}

func AggregatorValueCheck(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		return nil
	}
}

func SyncCommitteeValueCheck(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		return nil
	}
}

func SyncCommitteeContributionValueCheck(
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
