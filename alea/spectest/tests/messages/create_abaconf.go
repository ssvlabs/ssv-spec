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
		ExpectedRoot: "26fdaf9d6b419ed8e4feed76030ba8325d4210307d203dc4c103a75705620daa",
	}
}
