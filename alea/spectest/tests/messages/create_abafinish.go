package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
)

// CreateABAFinish tests creating a abafinish msg
func CreateABAFinish() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateABAFinish,
		Name:         "create abafinish",
		Vote:         byte(1),
		ExpectedRoot: "3b81faa20e3e985ddcef8133db776fad58f23e172cba08b849a7acc1758010c5",
	}
}
