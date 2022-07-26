package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateRoundChange,
		Name:         "create round change",
		Value:        []byte{1, 2, 3, 4},
		ExpectedRoot: "ce98c2909a7ce7e0372507f8c006862239f8c0382120037ca62f87a3947e5d20",
	}
}
