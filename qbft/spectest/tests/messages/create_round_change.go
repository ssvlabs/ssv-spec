package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() tests.SpecTest {
	return tests.NewCreateMsgSpecTest(
		"create round change",
		"Test creating a round change message when not previously prepared.",
		[32]byte{1, 2, 3, 4},
		nil,
		qbft.FirstRound,
		nil,
		nil,
		tests.CreateRoundChange,
		"a6ffc48674f1522fb90aa7bde2aa76cac54480cf366cdd4afcd7f8b4d548809a",
		nil,
		"",
	)
}
