package roundchange

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// JustificationMultiSigners tests a single prepare justification msg with multiple signers
func JustificationMultiSigners() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMultiSignerMessage(
			[]*rsa.PrivateKey{ks.NetworkKeys[2], ks.NetworkKeys[3]},
			[]types.OperatorID{2, 3},
		),
	}
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.NetworkKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "justification multi signer",
		Pre:            pre,
		PostRoot:       "ddb0f7e1a8888a8de5295005872d8525d6afea053121993dc65e34fcb7f290b2",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: round change justification invalid: msg allows 1 signer",
	}
}
