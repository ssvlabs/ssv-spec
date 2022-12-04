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

// HistoricalInstanceCapacity represents the upper bound of InstanceContainer a processmsg can process messages for as messages are not
// guaranteed to arrive in a timely fashion, we physically limit how far back the processmsg will process messages for
const HistoricalInstanceCapacity int = 5

type InstanceContainer [HistoricalInstanceCapacity]*Instance

func (i InstanceContainer) FindInstance(height Height) *Instance {
	for _, inst := range i {
		if inst != nil {
			if inst.GetHeight() == height {
				return inst
			}
		}
	}
	return nil
}

// addNewInstance will add the new instance at index 0, pushing all other stored InstanceContainer one index up (ejecting last one if existing)
func (i *InstanceContainer) addNewInstance(instance *Instance) {
	for idx := HistoricalInstanceCapacity - 1; idx > 0; idx-- {
		i[idx] = i[idx-1]
	}
	i[0] = instance
}
