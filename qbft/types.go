package qbft

import (
	"github.com/bloxapp/ssv-spec/p2p"
	"github.com/bloxapp/ssv-spec/types"
)

type Round uint64
type Height uint64

const (
	NoRound     Round  = 0 // NoRound represents a nil/ zero round
	FirstRound  Round  = 1 // FirstRound value is the first round in any QBFT instance start
	FirstHeight Height = 0
)

type Syncer interface {
	// SyncHighestDecided tries to fetch the highest decided from peers (not blocking)
	SyncHighestDecided(identifier types.MessageID) error
	// SyncHighestRoundChange tries to fetch for each committee member the highest round change broadcasted for the specific height from peers (not blocking)
	SyncHighestRoundChange(identifier types.MessageID, height Height) error
}

// Network is the interface for networking across QBFT components
type Network interface {
	Syncer
	p2p.Broadcaster
}

type Storage interface {
	// SaveHighestDecided saves (and potentially overrides) the highest Decided for a specific instance
	SaveHighestDecided(identifier types.MessageID, signedMsg *SignedMessage) error
	// GetHighestDecided returns highest decided if found, nil if didn't
	GetHighestDecided(identifier []byte) (*SignedMessage, error)
}

func ControllerIdToMessageID(identifier []byte) types.MessageID {
	ret := types.MessageID{}
	copy(ret[:], identifier)
	return ret
}
