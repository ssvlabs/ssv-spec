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
		CreateType:    tests.CreateVCBCAnswer,
		Name:          "create vcbcanswer",
		Proposals:     proposals,
		Priority:      alea.Priority(1),
		AggregatedMsg: []byte{5, 6, 7, 8},
		Author:        types.OperatorID(10),
		ExpectedRoot:  "a54e426ac86f2134101eaa30b722611e4069c790da72775c846f886c430c5206",
	}
}
