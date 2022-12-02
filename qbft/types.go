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

// Timer is an interface for a round timer, calling the UponRoundTimeout when times out
type Timer interface {
	// TimeoutForRound will reset running timer if exists and will start a new timer for a specific round
	TimeoutForRound(round Round)
}

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

func ControllerIdToMessageID(identifier []byte) types.MessageID {
	ret := types.MessageID{}
	copy(ret[:], identifier)
	return ret
}
