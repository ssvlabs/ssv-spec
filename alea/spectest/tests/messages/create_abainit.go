package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
)

// CreateABAInit tests creating a abainit msg
func CreateABAInit() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateABAInit,
		Name:         "create abainit",
		Vote:         byte(1),
		Round:        alea.Round(1),
		ExpectedRoot: "81cc6b475b941c46b07ba86362d8f049ed8bfd630ca648d8823e704736291be4",
	}
}
