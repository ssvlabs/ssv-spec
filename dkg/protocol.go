package dkg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type KeyGenOutcome struct {
	KeyGenOutput *KeyGenOutput
	BlameOutput  *BlameOutput
}

// KeyGenOutput is the bare minimum output from the protocol
type KeyGenOutput struct {
	Share           *bls.SecretKey
	OperatorPubKeys map[types.OperatorID]*bls.PublicKey
	ValidatorPK     types.ValidatorPK
	Threshold       uint64
}

// BlameOutput is the output of blame round
type BlameOutput struct {
	Valid        bool
	BlameMessage []byte
}

// KeyGenProtocol is an interface for all DKG protocol to support a variety of protocols for future upgrades
type KeyGenProtocol interface {
	Start(init *Init) error
	// ProcessMsg returns true and a bls share if finished
	ProcessMsg(msg *SignedMessage) (bool, *KeyGenOutcome, error)
}
