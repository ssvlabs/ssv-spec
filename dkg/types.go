package dkg

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"github.com/bloxapp/ssv-spec/types"
)

// Network is a collection of funcs for DKG
type Network interface {
	// StreamDKGOutput will stream to any subscriber the result of the DKG
	StreamDKGOutput(output *SignedOutput) error
	// Broadcast will broadcast a msg to the dkg network
	Broadcast(msg *SignedMessage) error
}

// Operator holds all info regarding a DKG Operator on the network
type Operator struct {
	// OperatorID the node's Operator ID
	OperatorID types.OperatorID
	// PubKey signing key for all message
	PubKey *ecdsa.PublicKey
	// EncryptionPubKey encryption pubkey for shares
	EncryptionPubKey *rsa.PublicKey
}

type Config struct {
	// Protocol the DKG protocol implementation
	Protocol            Protocol
	Network             Network
	SignatureDomainType types.DomainType
	Signer              types.DKGSigner
}
