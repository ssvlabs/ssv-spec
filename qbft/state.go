package qbft

import (
	"github.com/ssvlabs/ssv-spec/types"
)

type IConfig interface {
	// GetValueCheckF returns value check function
	GetValueCheckF() ProposedValueCheckF
	// GetProposerF returns func used to calculate proposer
	GetProposerF() ProposerF
	// GetNetwork returns a p2p Network instance
	GetNetwork() Network
	// GetTimer returns round timer
	GetTimer() Timer
	// GetCutOffRound returns the round that stops the instance
	GetCutOffRound() Round
}

type Config struct {
	Domain      types.DomainType
	ValueCheckF ProposedValueCheckF
	ProposerF   ProposerF
	Network     Network
	Timer       Timer
	CutOffRound Round
}

// GetSignatureDomainType returns the Domain type used for signatures
func (c *Config) GetSignatureDomainType() types.DomainType {
	return c.Domain
}

// GetValueCheckF returns value check instance
func (c *Config) GetValueCheckF() ProposedValueCheckF {
	return c.ValueCheckF
}

// GetProposerF returns func used to calculate proposer
func (c *Config) GetProposerF() ProposerF {
	return c.ProposerF
}

// GetNetwork returns a p2p Network instance
func (c *Config) GetNetwork() Network {
	return c.Network
}

// GetTimer returns round timer
func (c *Config) GetTimer() Timer {
	return c.Timer
}

func (c *Config) GetCutOffRound() Round {
	return c.CutOffRound
}

type State struct {
	CommitteeMember                 *types.CommitteeMember
	ID                              []byte // instance Identifier
	Round                           Round
	Height                          Height
	LastPreparedRound               Round
	LastPreparedValue               []byte
	ProposalAcceptedForCurrentRound *ProcessingMessage
	Decided                         bool
	DecidedValue                    []byte

	ProposeContainer     *MsgContainer
	PrepareContainer     *MsgContainer
	CommitContainer      *MsgContainer
	RoundChangeContainer *MsgContainer
}
