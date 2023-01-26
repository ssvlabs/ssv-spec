package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
)

// CreateABA tests creating a aba msg
func CreateABA() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateABA,
		Name:         "create aba",
		Vote:         byte(1),
		Round:        alea.Round(1),
		ExpectedRoot: "c6301e7506acd39fca837f0f64263b2e1eff17ccc683d87fa7293261c5eb50e1",
	}
}
