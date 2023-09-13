package types

import "github.com/attestantio/go-eth2-client/spec/phase0"

// DomainType is a unique identifier for signatures, 2 identical pieces of data signed with different domains will result in different sigs
type DomainType [4]byte

// Chains are intended to separate different SSV chains. A chain can have many forks in it.
const (
	MainnetChain = 0x0
	PrimusChain  = 0x1
	ShifuChain   = 0x2
	JatoChain    = 0x3
)

// DomainTypes represent specific forks for specific chains, messages are signed with the domain type making 2 messages from different domains incompatible
var (
	GenesisMainnet = DomainType{0x0, 0x0, MainnetChain, 0x0}
	PrimusTestnet  = DomainType{0x0, 0x0, PrimusChain, 0x0}
	ShifuTestnet   = DomainType{0x0, 0x0, ShifuChain, 0x0}
	ShifuV2Testnet = DomainType{0x0, 0x0, ShifuChain, 0x1}
	JatoTestnet    = DomainType{0x0, 0x0, JatoChain, 0x1}
)

// ForkData is a simple structure holding fork information for a specific chain (and its fork)
type ForkData struct {
	// Epoch in which the fork happened
	Epoch phase0.Epoch
	// Domain for the new fork
	Domain DomainType
}

func (domainType DomainType) GetChain() byte {
	return domainType[2]
}

func (domainType DomainType) GetForksData() []*ForkData {
	switch domainType.GetChain() {
	case MainnetChain:
		return mainnetForks()
	default:
		return []*ForkData{}
	}
}

// mainnetForks returns all forks for the mainnet chain
func mainnetForks() []*ForkData {
	return []*ForkData{
		{
			Epoch:  0,
			Domain: GenesisMainnet,
		},
	}
}
