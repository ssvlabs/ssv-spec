package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate rc msg (first one inserted, second not)
func DuplicateMsg() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 5),
		testingutils.TestingRoundChangeMessageWithParamsAndFullData(ks.OperatorKeys[1], types.OperatorID(1), 5, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs), testingutils.TestingQBFTFullData),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change duplicate msg",
		Pre:            pre,
		PostRoot:       "06ec71e4faa997e17856dc2af14e0494fc276bd85515ea45f833e238020e1d36",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
	}
}
