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
		ExpectedRoot: "c3cf296cf80a558fbc0a0de00759745f864268a123dc2e8e442f4883c54f0a82",
	}
}
