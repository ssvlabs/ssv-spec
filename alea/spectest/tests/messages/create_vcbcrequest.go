package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

// CreateVCBCRequest tests creating a vcbcrequest msg
func CreateVCBCRequest() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateVCBCRequest,
		Name:         "create vcbcrequest",
		Author:       types.OperatorID(10),
		Priority:     alea.Priority(2),
		ExpectedRoot: "3630543f3c58eaddb741f03ecc0ef57d2da67cd63f3c778c708b568f3a7fe8c7",
	}
}
