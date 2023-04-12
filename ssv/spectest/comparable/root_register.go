package ssvcomparable

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/ssv"
)

var RootRegister = map[string]ssv.Runner{}

// Register registers runner states and returns roots as strings
func Register(runners []ssv.Runner) []string {
	roots := make([]string, 0)
	for _, runner := range runners {
		r, err := runner.GetRoot()
		if err != nil {
			panic(err.Error())
		}
		roots = append(roots, hex.EncodeToString(r[:]))
		RootRegister[hex.EncodeToString(r[:])] = runner
	}
	return roots
}
