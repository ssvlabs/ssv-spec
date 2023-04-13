package qbftcomparable

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/qbft"
)

var RootRegister = map[string]*qbft.State{}

type StateComparison struct {
	PostState *qbft.State
}

// Register will register state roots with a global registry to be compared against
func (stateComp *StateComparison) Register() *StateComparison {
	r, err := stateComp.PostState.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	RootRegister[hex.EncodeToString(r[:])] = stateComp.PostState
	return stateComp
}

// Root returns all runner roots as string
func (stateComp *StateComparison) Root() string {
	if stateComp.PostState == nil {
		panic("state nil")
	}
	r, err := stateComp.PostState.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(r[:])
}
