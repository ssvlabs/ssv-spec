package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

// CreateFiller tests creating a filler msg
func CreateFiller() *tests.CreateMsgSpecTest {

	proposal1 := &alea.ProposalData{
		Data: []byte{1, 2, 3, 4},
	}
	proposal2 := &alea.ProposalData{
		Data: []byte{5, 6, 7, 8},
	}
	proposals1 := []*alea.ProposalData{proposal1, proposal2}

	proposal3 := &alea.ProposalData{
		Data: []byte{1, 2, 3, 4},
	}
	proposal4 := &alea.ProposalData{
		Data: []byte{5, 6, 7, 8},
	}
	proposals2 := []*alea.ProposalData{proposal3, proposal4}

	entries := [][]*alea.ProposalData{proposals1, proposals2}

	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateFiller,
		Name:         "create filler",
		Entries:      entries,
		Priorities:   []alea.Priority{alea.Priority(1), alea.Priority(2)},
		Author:       types.OperatorID(10),
		ExpectedRoot: "28d3649069510a18262882789bf4fc3dd37cb00052107ec4e650fbbd21a3cab2",
	}
}
