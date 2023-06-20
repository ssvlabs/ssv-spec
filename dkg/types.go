package dkg

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
)

// Network is a collection of funcs for DKG
type Network interface {
	// StreamDKGBlame will stream to any subscriber the blame result of the DKG
	StreamDKGBlame(blame *BlameOutput) error
	// StreamDKGOutput will stream to any subscriber the result of the DKG
	StreamDKGOutput(output map[types.OperatorID]*SignedOutput) error
	// BroadcastDKGMessage will broadcast a msg to the dkg network
	BroadcastDKGMessage(msg *SignedMessage) error
}

type Storage interface {
	// GetDKGOperator returns true and operator object if found by operator ID
	GetDKGOperator(operatorID types.OperatorID) (bool, *Operator, error)
	SaveKeyGenOutput(output *KeyGenOutput) error
	GetKeyGenOutput(pk types.ValidatorPK) (*KeyGenOutput, error)
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

type IConfig interface {
	// GetSigner returns a Signer instance
	GetSigner() types.DKGSigner
	// GetNetwork returns a p2p Network instance
	GetNetwork() Network
	// GetStorage returns a Storage instance
	GetStorage() Storage
}

type Config struct {
	// Protocol the DKG protocol implementation
	KeygenProtocol      func(RequestID, types.OperatorID, IConfig, *Init) Protocol
	ReshareProtocol     func(RequestID, types.OperatorID, IConfig, *Reshare, *ReshareParams) Protocol
	KeySign             func(RequestID, types.OperatorID, IConfig, *KeySign) Protocol
	Network             Network
	Storage             Storage
	SignatureDomainType types.DomainType
	Signer              types.DKGSigner
}

func (c *Config) GetSigner() types.DKGSigner {
	return c.Signer
}

// GetNetwork returns a p2p Network instance
func (c *Config) GetNetwork() Network {
	return c.Network
}

// GetStorage returns a Storage instance
func (c *Config) GetStorage() Storage {
	return c.Storage
}

type ErrInvalidRound struct{}

func (e ErrInvalidRound) Error() string {
	return "invalid dkg round"
}

type ErrMismatchRound struct{}

func (e ErrMismatchRound) Error() string {
	return "mismatch dkg round"
}
