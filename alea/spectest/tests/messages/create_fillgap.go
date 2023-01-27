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
		ExpectedRoot: "657e1c9f0774f2b611f0301bc4f6f9ebe40743169f361f4f134b67591f31d68d",
	}
}
