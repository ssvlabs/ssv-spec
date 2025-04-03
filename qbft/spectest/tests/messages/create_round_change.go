package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateRoundChange,
		Name:         "create round change",
		Value:        [32]byte{1, 2, 3, 4},
		ExpectedRoot: "05da5d07657444faa8e0de69b82cbcec63e8e52251887ab6cfff29da72f4adeb",
	}
}
