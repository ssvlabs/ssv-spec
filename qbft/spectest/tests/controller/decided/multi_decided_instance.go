package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiDecidedInstances tests deciding multiple instances
func MultiDecidedInstances() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()

	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue: []byte{1, 2, 3, 4},
			InputMessages: []*qbft.SignedMessage{
				testingutils.TestingCommitMultiSignerMessageWithHeight([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, height),
			},
			ExpectedDecidedState: tests.DecidedState{
				DecidedCnt: 1,
				DecidedVal: []byte{1, 2, 3, 4},
			},
			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "multi decide instances",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "c6e9b748e73de916edf6fb3a70c228f000ba014dc0a1bc54fa387fe528172fc8"),
			instanceData(1, "82e7de180912adfab8d5d5fbf1443b3bbedcfa16d1c35db4a1021e24db05e1d2"),
			instanceData(2, "db1af8f252886e69f88c66dd6e5c2c7ac790e539ff49b66aebe3f32dd4b7e8be"),
			instanceData(3, "cb70954bdae99971eac93a84f00f2fc792def306c71a86f9423f0234115fe0ed"),
			instanceData(4, "2c4d72bbaef85bcbbeede1e0cf05ad05c213b1ccf59d426351a196740ea8f3f6"),
			instanceData(5, "ca37672d0543116ec7fcd79ed3c15654acde54d0220b22cc40ba161a5f4b3099"),
			instanceData(6, "445521b9d6db9263821aeb0a4544a34d336fbc2c6d0561cdf6860602f3a73e15"),
			instanceData(7, "004c6b6a8dfbe9249115642af579a9e78877519eb950bd0192173987c1e99f8a"),
			instanceData(8, "2d500dd95118bd2d7b12b593acc126500e417198de9260c6e46eac117c85dc57"),
			instanceData(9, "2f89a0c05e5f04a00b65cbd9ba6adfe36342002bc76f935ab25384e4a0afedf3"),
			instanceData(10, "fed477fd0aa2f5bafd53acb422ae0d5a3571a308327acc8eefacc7d0d3ddb539"),
		},
	}
}
