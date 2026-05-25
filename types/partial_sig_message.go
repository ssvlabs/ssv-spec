package types

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type PartialSigMsgType uint64

const (
	// PostConsensusPartialSig is a partial signature over a decided duty (attestation data, block, etc)
	PostConsensusPartialSig = PartialSigMsgType(0)
	// RandaoPartialSig is a partial signature over randao reveal
	RandaoPartialSig = PartialSigMsgType(1)
	// ValidatorRegistrationPartialSig is a partial signature over a ValidatorRegistration object
	ValidatorRegistrationPartialSig = PartialSigMsgType(4)
	// VoluntaryExitPartialSig is a partial signature over a VoluntaryExit object
	VoluntaryExitPartialSig = PartialSigMsgType(5)
	// AggregatorCommitteePartialSig is a partial signature for combined aggregator and sync committee selection proofs
	AggregatorCommitteePartialSig = PartialSigMsgType(6)
)

type PartialSignatureMessages struct {
	Type PartialSigMsgType
	Slot phase0.Slot
	// 5048 = 3000 + 512 * 4 (worst-case scenario for
	// a committee with 3000 validators:
	// every validator has an aggregation duty
	// every validator is in sync committee and is a contributor to all 4 subnets
	Messages []*PartialSignatureMessage `ssz-max:"5048"`
}

// Encode returns a msg encoded bytes or error
func (msgs *PartialSignatureMessages) Encode() ([]byte, error) {
	return msgs.MarshalSSZ()
}

// Decode returns error if decoding failed
func (msgs *PartialSignatureMessages) Decode(data []byte) error {
	return msgs.UnmarshalSSZ(data)
}

// GetRoot returns the root used for signing and verification
func (msgs PartialSignatureMessages) GetRoot() ([32]byte, error) {
	return msgs.HashTreeRoot()
}

func (msgs PartialSignatureMessages) Validate() error {
	if len(msgs.Messages) == 0 {
		return NewError(NoPartialSigMessagesErrorCode, "no PartialSignatureMessages messages")
	}

	signer := msgs.Messages[0].Signer

	for _, m := range msgs.Messages {
		if signer != m.Signer {
			return NewError(InconsistentSignersErrorCode, "inconsistent signers")
		}
		if err := m.Validate(); err != nil {
			return errors.Wrap(err, "message invalid")
		}
	}
	return nil
}

// ValidateForSigner checks if the PartialSignatureMessages are valid for a given signer
// It checks if the signer is the same as the one in the messages
func (msgs PartialSignatureMessages) ValidateForSigner(signer OperatorID) error {
	if err := msgs.Validate(); err != nil {
		return err
	}
	if msgs.Messages[0].Signer != signer {
		return NewError(PartialSigInconsistentSignerErrorCode, "signer from signed message is inconsistent with partial signature signers")
	}
	return nil
}

// PartialSignatureMessage is a msg for partial Beacon chain related signatures (like partial attestation, block, randao sigs)
type PartialSignatureMessage struct {
	PartialSignature Signature `ssz-size:"96"` // The Beacon chain partial Signature for a duty
	SigningRoot      [32]byte  `ssz-size:"32"` // the root signed in PartialSignature
	Signer           OperatorID
	ValidatorIndex   phase0.ValidatorIndex
}

// Encode returns a msg encoded bytes or error
func (pcsm *PartialSignatureMessage) Encode() ([]byte, error) {
	return pcsm.MarshalSSZ()
}

// Decode returns error if decoding failed
func (pcsm *PartialSignatureMessage) Decode(data []byte) error {
	return pcsm.UnmarshalSSZ(data)
}

func (pcsm *PartialSignatureMessage) GetRoot() ([32]byte, error) {
	return pcsm.HashTreeRoot()
}

func (pcsm *PartialSignatureMessage) Validate() error {
	if pcsm.Signer == 0 {
		return NewError(ZeroSignerNotAllowedErrorCode, "signer ID 0 not allowed")
	}
	return nil
}
