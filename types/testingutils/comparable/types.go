package comparable

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/types"
)

type StateComparison struct {
	ExpectedState types.Root
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
