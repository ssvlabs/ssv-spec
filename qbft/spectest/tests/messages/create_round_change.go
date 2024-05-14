package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateRoundChange,
		Name:         "create round change",
		Value:        [32]byte{1, 2, 3, 4},
		ExpectedRoot: "da8a16c59e5c5a1e454ed318770e64d7e22b9209d1dbddb0445abca459cd5199",
	}
}
