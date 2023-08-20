package types

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type PartialSigMsgType uint64

const (
	// PostConsensusPartialSig is a partial signature over a decided duty (attestation data, block, etc)
	PostConsensusPartialSig PartialSigMsgType = iota
	// RandaoPartialSig is a partial signature over randao reveal
	RandaoPartialSig
	// SelectionProofPartialSig is a partial signature for aggregator selection proof
	SelectionProofPartialSig
	// ContributionProofs is the partial selection proofs for sync committee contributions (it's an array of sigs)
	ContributionProofs
	// ValidatorRegistrationPartialSig is a partial signature over a ValidatorRegistration object
	ValidatorRegistrationPartialSig
	// VoluntaryExitPartialSig is a partial signature over a VoluntaryExit object
	VoluntaryExitPartialSig
)

type PartialSignatureMessages struct {
	Type     PartialSigMsgType
	Slot     phase0.Slot
	Messages []*PartialSignatureMessage `ssz-max:"13"`
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
		return errors.New("no PartialSignatureMessages messages")
	}
	for _, m := range msgs.Messages {
		if err := m.Validate(); err != nil {
			return errors.Wrap(err, "message invalid")
		}
	}
	return nil
}

// PartialSignatureMessage is a msg for partial Beacon chain related signatures (like partial attestation, block, randao sigs)
type PartialSignatureMessage struct {
	PartialSignature Signature `ssz-size:"96"` // The Beacon chain partial Signature for a duty
	SigningRoot      [32]byte  `ssz-size:"32"` // the root signed in PartialSignature
	Signer           OperatorID
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
		return errors.New("signer ID 0 not allowed")
	}
	return nil
}

// SignedPartialSignatureMessage is an operator's signature over PartialSignatureMessage
type SignedPartialSignatureMessage struct {
	Message   PartialSignatureMessages
	Signature Signature `ssz-size:"96"`
	Signer    OperatorID
}

// Encode returns a msg encoded bytes or error
func (spcsm *SignedPartialSignatureMessage) Encode() ([]byte, error) {
	return spcsm.MarshalSSZ()
}

// Decode returns error if decoding failed
func (spcsm *SignedPartialSignatureMessage) Decode(data []byte) error {
	return spcsm.UnmarshalSSZ(data)
}

func (spcsm *SignedPartialSignatureMessage) GetSignature() Signature {
	return spcsm.Signature
}

func (spcsm *SignedPartialSignatureMessage) GetSigners() []OperatorID {
	return []OperatorID{spcsm.Signer}
}

func (spcsm *SignedPartialSignatureMessage) GetRoot() ([32]byte, error) {
	return spcsm.Message.GetRoot()
}

func (spcsm *SignedPartialSignatureMessage) Validate() error {
	if spcsm.Signer == 0 {
		return errors.New("signer ID 0 not allowed")
	}

	for _, msg := range spcsm.Message.Messages {
		if spcsm.Signer != msg.Signer {
			return errors.New("inconsistent signers")
		}
	}

	return spcsm.Message.Validate()
}
