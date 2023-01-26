package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

// CreateVCBCAnswer tests creating a vcbcanswer msg
func CreateVCBCAnswer() *tests.CreateMsgSpecTest {

	proposal1 := &alea.ProposalData{
		Data: []byte{1, 2, 3, 4},
	}
	proposal2 := &alea.ProposalData{
		Data: []byte{5, 6, 7, 8},
	}
	proposals := []*alea.ProposalData{proposal1, proposal2}

	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateVCBCAnswer,
		Name:         "create vcbcanswer",
		Proposals:    proposals,
		Priority:     alea.Priority(1),
		Proof:        types.Signature{},
		Author:       types.OperatorID(10),
		ExpectedRoot: "79724cc65ca79a26e4b97d989f654e9be37e2662a59f1e0b42c780fb13f1a1a2",
	}
}
