package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// EmptySigners tests a round change msg with no signers
func EmptySigners() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(5), 2)
	msg.OperatorIDs = []types.OperatorID{}

	msgs := []*types.SignedSSVMessage{
		msg,
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change empty signer",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: invalid SignedSSVMessage: no signers",
	}
}
