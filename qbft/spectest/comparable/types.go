package qbftcomparable

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/qbft"
)

// RootRegister is a global registry of expected state roots. e.g. {"0x123": &State{}}
var RootRegister = map[string]*qbft.RootGetter{}

type StateComparison struct {
	ExpectedState qbft.RootGetter
}

// Register will register state roots with a global registry to be compared against
func (stateComp *StateComparison) Register() *StateComparison {
	r, err := stateComp.ExpectedState.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	RootRegister[hex.EncodeToString(r[:])] = &stateComp.ExpectedState
	return stateComp
}

// Root returns all runner roots as string
func (stateComp *StateComparison) Root() string {
	if stateComp.ExpectedState == nil {
		panic("state nil")
	}
	r, err := stateComp.ExpectedState.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(r[:])
}
