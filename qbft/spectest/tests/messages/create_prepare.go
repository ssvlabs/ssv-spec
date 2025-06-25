package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// CreatePrepare tests creating a prepare msg
func CreatePrepare() tests.SpecTest {
	return tests.NewCreateMsgSpecTest(
		"create prepare",
		[32]byte{1, 2, 3, 4},
		nil,
		10,
		nil,
		nil,
		tests.CreatePrepare,
		"fe85e25b3cf7168e9e6417a9daaa71567a3c0c689b633d8154e252d8225c113c",
		nil,
		"",
	)
}
