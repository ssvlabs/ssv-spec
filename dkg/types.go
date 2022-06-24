package dkg

import (
	"crypto/rsa"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
)

// Network is a collection of funcs for DKG
type Network interface {
	// StreamDKGOutput will stream to any subscriber the result of the DKG
	StreamDKGOutput(output map[types.OperatorID]*SignedOutput) error
	// Broadcast will broadcast a msg to the dkg network
	Broadcast(msg types.Encoder) error

}

type Storage interface {
	// GetDKGOperator returns true and operator object if found by operator ID
	GetDKGOperator(operatorID types.OperatorID) (bool, *Operator, error)
}

// Operator holds all info regarding a DKG Operator on the network
type Operator struct {
	// OperatorID the node's Operator ID
	OperatorID types.OperatorID
	// ETHAddress the operator's eth address used to sign and verify messages against
	ETHAddress common.Address
	// EncryptionPubKey encryption pubkey for shares
	EncryptionPubKey *rsa.PublicKey
}

type ProtocolConfig struct {
	// Identifier unique for DKG session
	Identifier RequestID
	Operator *Operator
	BeaconNetwork       types.BeaconNetwork
	Signer              types.DKGSigner
}

type Config struct {
	// Protocol the DKG protocol implementation
	Protocol            func(init *Init, operatorID types.OperatorID, identifier RequestID) Protocol
	BeaconNetwork       types.BeaconNetwork
	Network             Network
	Storage             Storage
	SignatureDomainType types.DomainType
	Signer              types.DKGSigner
}
