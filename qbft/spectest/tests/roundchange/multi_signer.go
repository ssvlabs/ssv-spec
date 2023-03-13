package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiSigner tests a round change msg with multiple signers
func MultiSigner() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingMultiSignerRoundChangeMessageWithRound(
			[]*bls.SecretKey{ks.Shares[1], ks.Shares[2]},
			[]types.OperatorID{types.OperatorID(1), types.OperatorID(2)},
			2,
		),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change multi signers",
		Pre:            pre,
		PostRoot:       "ddb0f7e1a8888a8de5295005872d8525d6afea053121993dc65e34fcb7f290b2",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: msg allows 1 signer",
	}
}
