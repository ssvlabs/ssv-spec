package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CreateRoundChangePreviouslyPrepared tests creating a round change msg,previously prepared
func CreateRoundChangePreviouslyPrepared() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.CreateMsgSpecTest{
		CreateType: tests.CreateRoundChange,
		Name:       "create round change previously prepared",
		Round:      qbft.FirstRound,
		Value:      [32]byte{1, 2, 3, 4},
		PrepareJustifications: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
			testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
		},
		ExpectedRoot: "a6ffc48674f1522fb90aa7bde2aa76cac54480cf366cdd4afcd7f8b4d548809a",
	}
}
