package ssvcomparable

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/ssv"
)

type StateComparison struct {
	SyncCommitteeContribution ssv.Runner
	SyncCommittee             ssv.Runner
	Aggregator                ssv.Runner
	Proposer                  ssv.Runner
	BlindedProposer           ssv.Runner
	Attester                  ssv.Runner
	ValidatorRegistration     ssv.Runner
}

// Roots returns all runner roots as string
func (stateComp *StateComparison) Roots() map[ssv.Runner]string {
	roots := map[ssv.Runner]string{}
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

		roots[runner] = hex.EncodeToString(r[:])
	}
	return roots
}
