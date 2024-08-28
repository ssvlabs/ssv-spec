package ssv

import (
	"github.com/ssvlabs/ssv-spec/types"
)

type IConfig interface {
	// GetSignatureDomainType returns the Domain type used for signatures
	GetDomainType() types.DomainType
}

type Config struct {
	Domain types.DomainType
}

// GetDomainType returns the Domain type used for signatures
func (c *Config) GetDomainType() types.DomainType {
	return c.Domain
}
