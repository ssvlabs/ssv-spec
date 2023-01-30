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
		ExpectedRoot: "24748efbd7afdb3a46e978b9e7669c7b00d7df94b3f35b1aca2c5cc354259fac",
	}
}
