package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// InvalidValCheckData tests a decided message with invalid decided data (but should pass as it's decided)
func InvalidValCheckData() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide invalid value (should pass)",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*qbft.SignedMessage{
					// how to pass invalid value?
					testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 1,
					DecidedVal: testingutils.TestingInvalidValueCheck,
				},
				ControllerPostRoot: "befc8626cb71ca1b8493d41f54a709b22bbf530302177b8668fd77562ca3040e",
			},
		},
	}
}
