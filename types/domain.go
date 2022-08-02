package types

import "sync"

// DomainType is a unique identifier for signatures, 2 identical pieces of data signed with different domains will result in different sigs
type DomainType []byte

var (
	// PrimusTestnet is the domain for primus testnet
	PrimusTestnet = DomainType("primus_testnet")
	// ShifuTestnet is the domain for shifu testnet
	ShifuTestnet = DomainType("shifu")
)

var (
	domain DomainType
	once   sync.Once
)

// GetDefaultDomain returns the global domain
func GetDefaultDomain() DomainType {
	once.Do(func() {
		domain = ShifuTestnet
	})
	return domain
}

// SetDefaultDomain updates the global domain
func SetDefaultDomain(d DomainType) {
	domain = d
}
