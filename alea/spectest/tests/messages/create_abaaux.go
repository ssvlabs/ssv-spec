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
		ExpectedRoot: "f840d6122e04c1967f2f84235395eed6f058fc6192612f359ffc79639d52db21",
	}
}
