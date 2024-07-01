package roundchange

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NonUniqueSigner tests a round change msg with multiple signers and non unique signer
func NonUniqueSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingMultiSignerRoundChangeMessageWithRound(
		[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2]},
		[]types.OperatorID{types.OperatorID(1), types.OperatorID(2)},
		2,
	)
	msg.OperatorIDs = []types.OperatorID{types.OperatorID(1), types.OperatorID(1)}

	msgs := []*types.SignedSSVMessage{
		msg,
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change non unique signer",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: invalid SignedSSVMessage: non unique signer",
	}
}
