package qbftcomparable

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/qbft"
)

var RootRegister = map[string]*qbft.RootGetter{}

type StateComparison struct {
	qbft.RootGetter
}

// Register will register state roots with a global registry to be compared against
func (stateComp *StateComparison) Register() *StateComparison {
	r, err := stateComp.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	RootRegister[hex.EncodeToString(r[:])] = &stateComp.RootGetter

	return stateComp
}

// Root returns all runner roots as string
func (stateComp *StateComparison) Root() string {
	if stateComp.RootGetter == nil {
		panic("state nil")
	}
	r, err := stateComp.RootGetter.GetRoot()
	if err != nil {
		panic(err.Error())
	}

	return hex.EncodeToString(r[:])
}
