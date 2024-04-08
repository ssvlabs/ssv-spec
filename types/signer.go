package types

import (
	"bytes"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type SignatureDomain []byte
type Signature []byte

type SignatureType [4]byte

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
	SignBeaconObject(obj ssz.HashRoot, domain spec.Domain, pk []byte, domainType spec.DomainType) (Signature, [32]byte, error)
	// IsAttestationSlashable returns error if attestation is slashable
	IsAttestationSlashable(pk []byte, data *spec.AttestationData) error
	// IsBeaconBlockSlashable returns error if the given block is slashable
	IsBeaconBlockSlashable(pk []byte, slot spec.Slot) error
}

// OperatorSigner used for to sign protocol messages
type OperatorSigner interface {
	SignSSVMessage(data []byte, pk []byte) ([]byte, error)
}

// ShareSigner used for signing with the operator's share
type ShareSigner interface {
	SignRoot(data Root, sigType SignatureType, pk []byte) (Signature, error)
}

type DKGSigner interface {
	ShareSigner
	// SignDKGOutput signs output according to the SIP https://docs.google.com/document/d/1TRVUHjFyxINWW2H9FYLNL2pQoLy6gmvaI62KL_4cREQ/edit
	SignDKGOutput(output Root, address common.Address) (Signature, error)
	// SignETHDepositRoot signs an ethereum deposit root
	SignETHDepositRoot(root []byte, address common.Address) (Signature, error)
}

// KeyManager is an interface responsible for all key manager functions
type KeyManager interface {
	BeaconSigner
	ShareSigner
	// AddShare saves a share key
	AddShare(shareKey *bls.SecretKey) error
	// RemoveShare removes a share key
	RemoveShare(pubKey string) error
}
