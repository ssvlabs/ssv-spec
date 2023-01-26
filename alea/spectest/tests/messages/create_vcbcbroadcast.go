package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

// CreateVCBCBroadcast tests creating a vcbcbroadcast msg
func CreateVCBCBroadcast() *tests.CreateMsgSpecTest {

	proposal1 := &alea.ProposalData{
		Data: []byte{1, 2, 3, 4},
	}
	proposal2 := &alea.ProposalData{
		Data: []byte{5, 6, 7, 8},
	}
	proposals := []*alea.ProposalData{proposal1, proposal2}

	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateVCBCBroadcast,
		Name:         "create vcbcbroadcast",
		Proposals:    proposals,
		Priority:     alea.Priority(1),
		Author:       types.OperatorID(10),
		ExpectedRoot: "837ee1b4ac724afffee85d8155e419e9125539cbbde089639aa2a09393d19b91",
	}
}
