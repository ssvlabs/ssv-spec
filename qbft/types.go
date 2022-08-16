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

// Network is a collection of funcs for the QBFT Network
type Network interface {
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
