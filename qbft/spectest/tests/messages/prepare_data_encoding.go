package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PrepareDataEncoding tests encoding PrepareData
func PrepareDataEncoding() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1))

	r, _ := msg.GetRoot()
	b, _ := msg.Encode()

	return tests.NewMsgSpecTest(
		"prepare data encoding",
		[]*types.SignedSSVMessage{msg},
		[][]byte{b},
		[][32]byte{r},
		"",
	)
}
