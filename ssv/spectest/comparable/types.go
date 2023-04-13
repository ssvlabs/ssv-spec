package ssvcomparable

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/ssv"
)

var RootRegister = map[string]ssv.Runner{}

type StateComparison struct {
	SyncCommitteeContribution ssv.Runner
	SyncCommittee             ssv.Runner
	Aggregator                ssv.Runner
	Proposer                  ssv.Runner
	BlindedProposer           ssv.Runner
	Attester                  ssv.Runner
	ValidatorRegistration     ssv.Runner
}

// Register will register state roots with a global registry to be compared against
func (stateComp *StateComparison) Register() *StateComparison {
	for _, runner := range []ssv.Runner{
		stateComp.SyncCommitteeContribution,
		stateComp.SyncCommittee,
		stateComp.Aggregator,
		stateComp.Proposer,
		stateComp.BlindedProposer,
		stateComp.Attester,
		stateComp.ValidatorRegistration,
	} {
		if runner != nil {
			stateComp.register(runner)
		}
	}
	return stateComp
}

// Roots returns all runner roots as string
func (stateComp *StateComparison) Roots() []string {
	roots := make([]string, 0)
	for _, runner := range []ssv.Runner{
		stateComp.SyncCommitteeContribution,
		stateComp.SyncCommittee,
		stateComp.Aggregator,
		stateComp.Proposer,
		stateComp.BlindedProposer,
		stateComp.Attester,
		stateComp.ValidatorRegistration,
	} {
		if runner == nil {
			continue
		}
		r, err := runner.GetRoot()
		if err != nil {
			panic(err.Error())
		}
		roots = append(roots, hex.EncodeToString(r[:]))
	}
	return roots
}

// register runner states
func (stateComp *StateComparison) register(runner ssv.Runner) {
	r, err := runner.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	RootRegister[hex.EncodeToString(r[:])] = runner
}
