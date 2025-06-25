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

	inputMessages := []*types.SignedSSVMessage{
		msg,
	}

	return tests.NewMsgProcessingSpecTest(
		"round change non unique signer",
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: invalid SignedSSVMessage: non unique signer",
		nil,
	)
}
