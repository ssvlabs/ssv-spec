package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

// CreateVCBCReady tests creating a vcbcready msg
func CreateVCBCReady() *tests.CreateMsgSpecTest {

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
		CreateType:   tests.CreateVCBCReady,
		Name:         "create vcbcready",
		Hash:         hash,
		Priority:     alea.Priority(1),
		Author:       types.OperatorID(10),
		ExpectedRoot: "bfab94740b0565179fa0d89b62435c67bf8c98f6e78906351d646727df1e2af4",
	}
}
