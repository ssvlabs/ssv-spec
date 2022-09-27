package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
)

type Round uint64
type Height int64

const (
	NoRound     Round  = 0 // NoRound represents a nil/ zero round
	FirstRound  Round  = 1 // FirstRound value is the first round in any QBFT instance start
	FirstHeight Height = 0
)

type Sync interface {
	// SyncHighestDecided tries to fetch the highest decided from peers (not blocking)
	SyncHighestDecided(identifier []byte) error
	// SyncHighestRoundChange tries to fetch for each committee member the highest round change broadcasted for the specific height from peers (not blocking)
	SyncHighestRoundChange(identifier []byte, height Height) error
}

// Network is a collection of funcs for the QBFT Network
type Network interface {
	Sync
	Broadcast(msg types.Encoder) error
	BroadcastDecided(msg types.Encoder) error
}

type Storage interface {
	// SaveHighestDecided saves (and potentially overrides) the highest Decided for a specific instance
	SaveHighestDecided(signedMsg *SignedMessage) error
	// GetHighestDecided returns highest decided if found, nil if didn't
	GetHighestDecided(identifier []byte) (*SignedMessage, error)
}

func ControllerIdToMessageID(identifier []byte) types.MessageID {
	ret := types.MessageID{}
	copy(ret[:], identifier)
	return ret
}
