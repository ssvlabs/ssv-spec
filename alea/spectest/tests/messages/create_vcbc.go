package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
)

// CreateVCBC tests creating a vcbc msg
func CreateVCBC() *tests.CreateMsgSpecTest {

	proposal1 := &alea.ProposalData{
		Data: []byte{1, 2, 3, 4},
	}
	proposal2 := &alea.ProposalData{
		Data: []byte{5, 6, 7, 8},
	}
	proposals := []*alea.ProposalData{proposal1, proposal2}
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateVCBC,
		Name:         "create vcbc",
		Proposals:    proposals,
		Priority:     0,
		ExpectedRoot: "b696cbf25c3dea9cac3de29d31f9aecca8bf7734002e3ba0f5574b066da42990",
	}
}
