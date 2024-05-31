package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownSigner tests a signed round change msg with an unknown signer
func UnknownSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(5), 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change unknown signer",
		Pre:            pre,
		PostRoot:       "50829de8e9d4814e529e39647a5faf1dd7c93dad00e3b59ed1215559b408671f",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: signer not in committee",
	}
}
