package keysign

import (
	"context"
	"sync"

	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/bloxapp/ssv-spec/dkg/frost"
)

// State tracks protocol's current round, stores messages in MsgContainer, stores
// session key and operator's secret shares
type State struct {
	// round mutex ensures atomic access to current round
	roundMutex   *sync.Mutex
	currentRound common.ProtocolRound

	// a container to store messages for each round from each operator
	msgContainer common.IMsgContainer

	// underlying timer for timeout
	roundTimer *frost.RoundTimer
}

func initState() *State {
	return &State{
		currentRound: common.Uninitialized,
		msgContainer: common.NewMsgContainer(),
		roundTimer:   frost.NewRoundTimer(context.Background(), nil),
		roundMutex:   new(sync.Mutex),
	}
}

func (state *State) GetCurrentRound() common.ProtocolRound {
	state.roundMutex.Lock()
	defer state.roundMutex.Unlock()

	return state.currentRound
}

func (state *State) SetCurrentRound(round common.ProtocolRound) {
	state.roundMutex.Lock()
	defer state.roundMutex.Unlock()

	state.currentRound = round
}
