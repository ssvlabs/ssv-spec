package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

// CreateVCBCFinal tests creating a vcbcfinal msg
func CreateVCBCFinal() *tests.CreateMsgSpecTest {

	proposal1 := &alea.ProposalData{
		Data: []byte{1, 2, 3, 4},
	}
	proposal2 := &alea.ProposalData{
		Data: []byte{5, 6, 7, 8},
	}
	proposals := []*alea.ProposalData{proposal1, proposal2}

	hash, err := alea.GetProposalsHash(proposals)
	if err != nil {
		errors.Wrap(err, "could not generate hash from proposals")
	}

	return &tests.CreateMsgSpecTest{
		CreateType:    tests.CreateVCBCFinal,
		Name:          "create vcbcfinal",
		Hash:          hash,
		Priority:      alea.Priority(1),
		AggregatedMsg: []byte{5, 6, 7, 8},
		Author:        types.OperatorID(10),
		ExpectedRoot:  "c3b32c63bcffbfbb788ed58c937dda7cce20131246bb0059cfc10b0b673aeae4",
	}
}
