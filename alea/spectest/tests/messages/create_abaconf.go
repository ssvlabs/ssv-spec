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
		ExpectedRoot: "fb876db4888fd020bbec1d715d2ed813872bc550b45ad81b68fcc64372d8bf1b",
	}
}
