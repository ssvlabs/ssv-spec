package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
)

// CreateABAConf tests creating a abaconf msg
func CreateABAConf() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateABAConf,
		Name:         "create abaconf",
		Votes:        []byte{0, 1},
		Round:        alea.Round(1),
		ExpectedRoot: "5dbe678d18f27700ffc95909614f227fda250130e018a28c103b5065eaaaaeea",
	}
}
