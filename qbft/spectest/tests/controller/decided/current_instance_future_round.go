package decided

import (
	"crypto/rsa"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CurrentInstanceFutureRound tests a decided msg received for current running instance for a future round
func CurrentInstanceFutureRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide current instance future round",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),
					testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
					testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
					testingutils.TestingPrepareMessage(ks.OperatorKeys[3], 3),
					testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
					testingutils.TestingCommitMessage(ks.OperatorKeys[2], 2),
					testingutils.TestingCommitMultiSignerMessageWithRound([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}, 50),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 1,
					DecidedVal: testingutils.TestingQBFTFullData,
				},
			},
		},
	}
}
