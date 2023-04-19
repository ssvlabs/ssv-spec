package qbftcomparable

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/qbft"
)

var RootRegister = map[string]*StateComparison{}

type StateComparison struct {
	PostState      *qbft.State
	PostController *qbft.Controller
}

// Register will register state roots with a global registry to be compared against
func (stateComp *StateComparison) Register() *StateComparison {
	var r []byte
	var err error

	if stateComp.PostState != nil {
		r, err = stateComp.PostState.GetRoot()
		if err != nil {
			panic(err.Error())
		}
	}

	if stateComp.PostController != nil {
		r, err = stateComp.PostController.GetRoot()
		if err != nil {
			panic(err.Error())
		}
	}

	if r == nil {
		panic("nil root")
	}

	RootRegister[hex.EncodeToString(r[:])] = stateComp
	return stateComp
}

// Root returns all runner roots as string
func (stateComp *StateComparison) Root() string {
	var r []byte
	var err error

	if stateComp.PostState == nil && stateComp.PostController == nil {
		panic("state nil")
	}

	if stateComp.PostState != nil {
		r, err = stateComp.PostState.GetRoot()
	}

	if stateComp.PostController != nil {
		r, err = stateComp.PostController.GetRoot()
	}

	if err != nil {
		panic(err.Error())
	}

	return hex.EncodeToString(r[:])
}
