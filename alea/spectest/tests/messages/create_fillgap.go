package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

// CreateFillGap tests creating a fillgap msg
func CreateFillGap() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateFillGap,
		Name:         "create fillgap",
		Author:       types.OperatorID(10),
		Priority:     alea.Priority(2),
		ExpectedRoot: "05883fb7b1c121c5edd97f827a1d71468376b4dd713d8720c1b4f977cf249f4b",
	}
}
