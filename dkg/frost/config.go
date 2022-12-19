package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
)

type IConfig interface {
	// GetSigner returns a Signer instance
	GetSigner() types.DKGSigner
	// GetNetwork returns a p2p Network instance
	GetNetwork() dkg.Network
	// GetStorage returns a Storage instance
	GetStorage() dkg.Storage
}

type Config struct {
	network dkg.Network
	signer  types.DKGSigner
	storage dkg.Storage
}

// GetSigner returns a Signer instance
func (c *Config) GetSigner() types.DKGSigner {
	return c.signer
}

// GetNetwork returns a p2p Network instance
func (c *Config) GetNetwork() dkg.Network {
	return c.network
}

// GetStorage returns a Storage instance
func (c *Config) GetStorage() dkg.Storage {
	return c.storage
}
