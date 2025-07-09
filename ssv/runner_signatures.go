package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

func (b *BaseRunner) signBeaconObject(runner Runner, duty *types.ValidatorDuty,
	obj ssz.HashRoot, slot spec.Slot, domainType spec.DomainType) (*types.PartialSignatureMessage, error) {
	epoch := runner.GetBaseRunner().BeaconNetwork.EstimatedEpochAtSlot(slot)
	domain, err := runner.GetBeaconNode().DomainData(epoch, domainType)
	if err != nil {
		return nil, errors.Wrap(err, "could not get beacon domain")
	}

	sig, r, err := runner.GetSigner().SignBeaconObject(obj, domain,
		runner.GetBaseRunner().Share[duty.ValidatorIndex].SharePubKey,
		domainType)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign beacon object")
	}

	return &types.PartialSignatureMessage{
		PartialSignature: sig,
		SigningRoot:      r,
		Signer:           runner.GetOperatorSigner().GetOperatorID(),
		ValidatorIndex:   duty.ValidatorIndex,
	}, nil
}

// Validate message content without verifying signatures
func (b *BaseRunner) validatePartialSigMsgForSlot(
	psigMsgs *types.PartialSignatureMessages,
	slot spec.Slot,
) error {
	if err := psigMsgs.Validate(); err != nil {
		return errors.Wrap(err, "PartialSignatureMessages invalid")
	}
	if psigMsgs.Slot != slot {
		return errors.New("invalid partial sig slot")
	}

	// Get signer
	msgSigner := psigMsgs.Messages[0].Signer // signer is the same in all psigMsgs.Messages and len(psigMsgs.Messages) > 0 (guaranteed by psigMsgs.Validate())

	// Get committee (unique for runner)
	var shareSample *types.Share
	for _, share := range b.Share {
		shareSample = share
		break
	}
	if shareSample == nil {
		return errors.New("can not get committee because there is no share in runner")
	}
	committee := shareSample.Committee

	// Check if signer is in committee
	signerInCommittee := false
	for _, operator := range committee {
		if operator.Signer == msgSigner {
			signerInCommittee = true
			break
		}
	}
	if !signerInCommittee {
		return errors.New("unknown signer")
	}

	return nil
}

// Validate if runner has a share for each ValidatorIndex in the PartialSignatureMessages.Messages
func (b *BaseRunner) validateValidatorIndexInPartialSigMsg(
	psigMsgs *types.PartialSignatureMessages,
) error {
	for _, msg := range psigMsgs.Messages {
		// Check if it has the validator index share
		_, ok := b.Share[msg.ValidatorIndex]
		if !ok {
			return errors.New("unknown validator index")
		}
	}
	return nil
}

func (b *BaseRunner) verifyBeaconPartialSignature(signer types.OperatorID, signature types.Signature, root [32]byte,
	committee []*types.ShareMember) error {
	for _, n := range committee {
		if n.Signer == signer {
			pk := &bls.PublicKey{}
			if err := pk.Deserialize(n.SharePubKey); err != nil {
				return errors.Wrap(err, "could not deserialized pk")
			}
			sig := &bls.Sign{}
			if err := sig.Deserialize(signature); err != nil {
				return errors.Wrap(err, "could not deserialized Signature")
			}

			// verify
			if !sig.VerifyByte(pk, root[:]) {
				return errors.New("wrong signature")
			}
			return nil
		}
	}
	return errors.New("unknown signer")
}

// Stores the container's existing signature or the new one, depending on their validity. If both are invalid, remove the existing one
func (b *BaseRunner) resolveDuplicateSignature(container *PartialSigContainer, msg *types.PartialSignatureMessage) {
	// Check previous signature validity
	previousSignature, err := container.GetSignature(msg.ValidatorIndex, msg.Signer, msg.SigningRoot)
	if err == nil {
		err = b.verifyBeaconPartialSignature(msg.Signer, previousSignature, msg.SigningRoot,
			b.Share[msg.ValidatorIndex].Committee)
		if err == nil {
			// Keep the previous sigature since it's correct
			return
		}
	}

	// Previous signature is incorrect or doesn't exist
	container.Remove(msg.ValidatorIndex, msg.Signer, msg.SigningRoot)

	// Hold the new signature, if correct
	err = b.verifyBeaconPartialSignature(msg.Signer, msg.PartialSignature, msg.SigningRoot,
		b.Share[msg.ValidatorIndex].Committee)
	if err == nil {
		container.AddSignature(msg)
	}
}
