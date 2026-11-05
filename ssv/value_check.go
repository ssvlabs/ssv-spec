package ssv

import (
	"bytes"
	"fmt"

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

// BeaconVoteValueCheckF validates the committee QBFT value (SIP #94 §2): the vote's invariants, the expected
// checkpoints, and slashability. The slashability data carries the decided attestation index, so a vote that
// differs from a signed one only in the index is a double vote.
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
			Slot:            slot,
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

// ProposerValueCheckF validates the proposer QBFT value. runningDutySlot reports the running duty's slot so a
// value for any other slot is rejected before consensus can commit it. It is the only slot guard (there is no
// post-decide backstop), so production callers must pass a real provider; nil skips the bind and is for
// isolated value-check tests only.
func ProposerValueCheckF(
	signer types.BeaconSigner,
	network types.BeaconNetwork,
	validatorPK types.ValidatorPK,
	validatorIndex phase0.ValidatorIndex,
	sharePublicKey []byte,
	runningDutySlot func() phase0.Slot,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		cd := &types.ProposerConsensusData{}
		if err := cd.Decode(data); err != nil {
			return types.WrapError(types.ProposerConsensusDataDecodeErrorCode, errors.Wrap(err, "failed decoding consensus data"))
		}
		// Bind the value to the running duty's slot (SIP #94 §4): QBFT decides whatever the leader proposes, and
		// an instance decided on another slot's value can never serve the running duty. A nil provider (isolated
		// value-check tests) or a 0 slot (no running duty) skips the bind.
		if runningDutySlot != nil {
			if want := runningDutySlot(); want != 0 && cd.Duty.Slot != want {
				return types.NewError(types.ProposerDutySlotMismatchErrorCode, "consensus data duty slot does not match running duty slot")
			}
		}

		// Validate pins the leader-stamped Version to Gloas (SIP #94 §4) and decodes the proposal.
		if err := cd.Validate(); err != nil {
			return types.NewError(types.QBFTValueInvalidErrorCode, fmt.Sprintf("invalid value: %v", err.Error()))
		}

		if err := dutyValueCheck(&cd.Duty, network, types.BNRoleProposer, validatorPK, validatorIndex); err != nil {
			return errors.Wrap(err, "duty invalid")
		}

		proposalData, err := cd.GetBlockData()
		if err != nil {
			return errors.Wrap(err, "could not get block data")
		}
		block := proposalData.Block
		// The block must be for the duty's slot (SIP #94 §4).
		if block.Slot != cd.Duty.Slot {
			return types.NewError(types.ProposerBlockSlotMismatchErrorCode, "block slot does not match duty slot")
		}
		// The block must name the duty's validator, whose key signs it; the beacon node would reject any other
		// proposer index, losing the slot (SIP #94 §4).
		if block.ProposerIndex != cd.Duty.ValidatorIndex {
			return types.NewError(types.ProposerBlockProposerIndexMismatchErrorCode, "block proposer index does not match duty validator index")
		}
		// payload_root is non-zero iff the bid is self-build (SIP #94 §4).
		selfBuild := block.Body.SignedExecutionPayloadBid.Message.BuilderIndex == gloas.BuilderIndexSelfBuild
		payloadZero := proposalData.PayloadRoot == phase0.Root{}
		if selfBuild == payloadZero {
			return types.NewError(types.QBFTValueInvalidErrorCode, "payload_root presence does not match self-build bid")
		}
		return signer.IsBeaconBlockSlashable(sharePublicKey, block.Slot)
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
