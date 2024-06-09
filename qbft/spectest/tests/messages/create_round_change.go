package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateRoundChange,
		Name:         "create round change",
		Value:        [32]byte{1, 2, 3, 4},
		ExpectedRoot: "78b5a214f4b2c83afe7dc739b317327a6fe8a29e14c27293fbd8a5522543dbf7",
	}
}
