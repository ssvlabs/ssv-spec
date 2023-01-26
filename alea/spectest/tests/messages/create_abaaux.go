package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
)

// CreateABAAux tests creating a abaaux msg
func CreateABAAux() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateABAAux,
		Name:         "create abaaux",
		Vote:         byte(1),
		Round:        alea.Round(1),
		ExpectedRoot: "81cc6b475b941c46b07ba86362d8f049ed8bfd630ca648d8823e704736291be4",
	}
}
