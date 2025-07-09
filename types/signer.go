package types

import (
	"bytes"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

type SignatureDomain []byte
type Signature []byte
type SignatureType [4]byte
type SignSSVMessageF = func(ssvMessage *SSVMessage) ([]byte, error)

func (sigType SignatureType) Equal(other SignatureType) bool {
	return bytes.Equal(sigType[:], other[:])
}

var (
	QBFTSignatureType    SignatureType = [4]byte{1, 0, 0, 0}
	PartialSignatureType SignatureType = [4]byte{2, 0, 0, 0}
	DKGSignatureType     SignatureType = [4]byte{3, 0, 0, 0}
)

type BeaconSigner interface {
	// SignBeaconObject returns signature and root.
	SignBeaconObject(obj ssz.HashRoot, domain phase0.Domain, pk []byte, domainType DomainType) (Signature, [32]byte, error)
	// IsAttestationSlashable returns error if attestation is slashable
	IsAttestationSlashable(pk ShareValidatorPK, data *phase0.AttestationData) error
	// IsBeaconBlockSlashable returns error if the given block is slashable
	IsBeaconBlockSlashable(pk []byte, slot phase0.Slot) error
}
