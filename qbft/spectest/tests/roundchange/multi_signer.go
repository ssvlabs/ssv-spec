package roundchange

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MultiSigner tests a round change msg with multiple signers
func MultiSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingMultiSignerRoundChangeMessageWithRound(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2]},
			[]types.OperatorID{types.OperatorID(1), types.OperatorID(2)},
			2,
		),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change multi signers",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: msg allows 1 signer",
	}
}
