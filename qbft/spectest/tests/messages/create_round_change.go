package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateRoundChange,
		Name:         "create round change",
		Round:        qbft.FirstRound,
		Value:        [32]byte{1, 2, 3, 4},
		ExpectedRoot: "a6ffc48674f1522fb90aa7bde2aa76cac54480cf366cdd4afcd7f8b4d548809a",
	}
}
