package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// CurrentInstance tests a decided msg received for current running instance
func CurrentInstance() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.ControllerSpecTest{
		Name: "decide current instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingProposalMessage(ks.Shares[1], 1),

					testingutils.TestingPrepareMessage(ks.Shares[1], 1),
					testingutils.TestingPrepareMessage(ks.Shares[2], 2),
					testingutils.TestingPrepareMessage(ks.Shares[3], 3),

					testingutils.TestingCommitMessage(ks.Shares[1], 1),
					testingutils.TestingCommitMessage(ks.Shares[2], 2),
					testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 1,
					DecidedVal: testingutils.TestingQBFTFullData,
				},
				ControllerPostRoot: "da1e383ad3be736cb63fe62d0981576dacb99166530c008292f0718852d77ffb",
			},
		},
	}
}
