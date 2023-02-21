package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CreateRoundChangePreviouslyPrepared tests creating a round change msg,previously prepared
func CreateRoundChangePreviouslyPrepared() *tests.CreateMsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.CreateMsgSpecTest{
		CreateType: tests.CreateProposal,
		Name:       "create round change previously prepared",
		Value:      [32]byte{1, 2, 3, 4},
		PrepareJustifications: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
			testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
		},
		ExpectedRoot: "de421aca0c42404cea60b6a9dea40457831eed611b78dd529e022fc54e81af44",
	}
}
