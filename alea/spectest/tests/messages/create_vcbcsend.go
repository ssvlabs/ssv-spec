package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

// CreateVCBCSend tests creating a vcbcsend msg
func CreateVCBCSend() *tests.CreateMsgSpecTest {

	proposal1 := &alea.ProposalData{
		Data: []byte{1, 2, 3, 4},
	}
	proposal2 := &alea.ProposalData{
		Data: []byte{5, 6, 7, 8},
	}
	proposals := []*alea.ProposalData{proposal1, proposal2}
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateVCBCSend,
		Name:         "create vcbcsend",
		Proposals:    proposals,
		Priority:     alea.Priority(1),
		Author:       types.OperatorID(10),
		ExpectedRoot: "dfe94cdfffc195b06de569e962a1dee41488db83a7910f274473d6fe8bdb3e3d",
	}
}
