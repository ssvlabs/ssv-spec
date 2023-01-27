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
		ExpectedRoot: "eb24ad0d50549b80ff747ba2e7fdfbfbbd4029de6b18580889f0857780dc2bcb",
	}
}
