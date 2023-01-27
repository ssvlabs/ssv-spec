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
		ExpectedRoot: "05199530233dfcc27804128c18f3f0cad9e5f06cf7041fbeb492857c6bdb7c91",
	}
}
