package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// JustificationNoQuorum tests justifications with no quorum
func JustificationNoQuorum() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
	}
	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change justification no quorum",
		Pre:            pre,
		PostRoot:       "ddb0f7e1a8888a8de5295005872d8525d6afea053121993dc65e34fcb7f290b2",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: no justifications quorum",
	}
}
