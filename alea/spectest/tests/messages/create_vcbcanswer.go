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
		ExpectedRoot: "cec515cc38c1905933214aedae8653963c5e27ea224c49ce05ab76ab4128eb9b",
	}
}
