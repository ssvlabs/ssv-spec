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
		CreateType:   tests.CreateVCBCFinal,
		Name:         "create vcbcfinal",
		Hash:         hash,
		Priority:     alea.Priority(1),
		Proof:        types.Signature{},
		Author:       types.OperatorID(10),
		ExpectedRoot: "419e91fbd4825dd4dd1ed0b0305dd4ccfc08bc452ac7f6498ebd5965d3ad707b",
	}
}
