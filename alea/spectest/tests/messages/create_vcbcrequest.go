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
		ExpectedRoot: "cffcc7bbaabb55f6f9d71d75d13b8f8ba824d2b7b838dab8639de2ea19628c86",
	}
}
