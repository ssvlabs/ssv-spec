package hbbft

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/p2p"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

type Round uint64
type Height uint64
type ACSRound uint64

const (
	NoRound           Round    = 0 // NoRound represents a nil/ zero round
	FirstRound        Round    = 1 // FirstRound value is the first round in any QBFT instance start
	FirstACSRound     ACSRound = 1
	FirstHeight       Height   = 0
	HBBFTDefaultRound Round    = 1
)

// Timer is an interface for a round timer, calling the UponRoundTimeout when times out
type Timer interface {
	// TimeoutForRound will reset running timer if exists and will start a new timer for a specific round
	TimeoutForRound(round Round)
}

type Syncer interface {
	// SyncHighestDecided tries to fetch the highest decided from peers (not blocking)
	SyncHighestDecided(identifier types.MessageID) error
	// SyncDecidedByRange will trigger sync from-to heights (including)
	SyncDecidedByRange(identifier types.MessageID, to, from Height)
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

type InstanceContainer []*Instance

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
	indexToInsert := len(*i)
	for index, existingInstance := range *i {
		if existingInstance.GetHeight() < instance.GetHeight() {
			indexToInsert = index
			break
		}
	}
	*i = insertAtIndex(*i, indexToInsert, instance)
}

func insertAtIndex(arr []*Instance, index int, value *Instance) InstanceContainer {
	if len(arr) == index { // nil or empty slice or after last element
		return append(arr, value)
	}
	arr = append(arr[:index+1], arr[index:]...) // index < len(a)
	arr[index] = value
	return arr
}
