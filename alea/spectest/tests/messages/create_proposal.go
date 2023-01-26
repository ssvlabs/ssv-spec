package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
)

// CreateProposal tests creating a proposal msg
func CreateProposal() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateProposal,
		Name:         "create proposal",
		Value:        []byte{1, 2, 3, 4},
		ExpectedRoot: "7ba23d5d30cce422bf99147c458942000e7dac1ae3539ad77d2333803e4037f6",
	}
}
