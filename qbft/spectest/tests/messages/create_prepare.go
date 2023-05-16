package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreatePrepare tests creating a prepare msg
func CreatePrepare() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:      tests.CreatePrepare,
		Name:            "create prepare",
		Value:           []byte{1, 2, 3, 4},
		Round:           10,
		ExpectedSSZRoot: "f1c182dc944840b3cdd182f82ce891a9b527cade2033ff5b371e1ee1086841e2",
	}
}
