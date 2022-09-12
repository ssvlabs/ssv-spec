package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateRoundChange,
		Name:         "create round change",
		Value:        []byte{1, 2, 3, 4},
		ExpectedRoot: "6bc0f189b080c79e6cadea6a2de13f7b79b38589af73297c794742e1c8c9b13c",
	}
}
