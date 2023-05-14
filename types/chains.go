package types

import "github.com/attestantio/go-eth2-client/spec/phase0"

// DomainType is a unique identifier for signatures, 2 identical pieces of data signed with different domains will result in different sigs
type DomainType [4]byte

const (
	GenesisChain = 0x0
	PrimusChain  = 0x1
	ShifuChain   = 0x2
	JatoChain    = 0x3
)

var (
	GenesisMainnet = DomainType{0x0, 0x0, GenesisChain, 0x0}
	PrimusTestnet  = DomainType{0x0, 0x0, PrimusChain, 0x0}
	ShifuTestnet   = DomainType{0x0, 0x0, ShifuChain, 0x0}
	ShifuV2Testnet = DomainType{0x0, 0x0, ShifuChain, 0x1}
	JatoTestnet    = DomainType{0x0, 0x0, JatoChain, 0x1}
)

type ForkData struct {
	Epoch  phase0.Epoch
	Domain DomainType
}

func (domainType DomainType) GetChain() byte {
	return domainType[2]
}

func (domainType DomainType) GetForksData() []*ForkData {
	switch domainType.GetChain() {
	case GenesisChain:
		return genesisForks()
	default:
		return []*ForkData{}
	}
}

func genesisForks() []*ForkData {
	return []*ForkData{
		{
			Epoch:  0,
			Domain: GenesisMainnet,
		},
	}
}
