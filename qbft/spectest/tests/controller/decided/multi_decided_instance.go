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
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue: []byte{1, 2, 3, 4},
			InputMessages: []*qbft.SignedMessage{
				testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     height,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
					}),
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
			instanceData(qbft.FirstHeight, "73ba8a44f10c67c1885385c76076fea5b57c2561c0d506f6d15ec62414f38591"),
			instanceData(1, "d233358de3510f6236f9441d8cf14a97b5e3b7ce5bdf410090ef660d74a7583a"),
			instanceData(2, "0ca69a226ad6acd5f2867172b00e5f69b1935c7986e2f3981137b255ce44e8dd"),
			instanceData(3, "9da8c7102a4b7d66df477d8dbd22b2665eb2891fafc1482c67d491dfe7406120"),
			instanceData(4, "c75b8981c604389e191c4b36a6db6537b1becba91a65df1e4af1a237475fb7f3"),
			instanceData(5, "ca0d84be5214cb15e871edd3885ba03153b3420ffabc13babf9a61b6af0989aa"),
			instanceData(6, "f8dac0bc1ce17dd65faed261553fb8a30ab78e630807dfc2696ec85088b8df0a"),
			instanceData(7, "60d538a3b3b03cf2b9e1a69bd8a9ea0fa4f05f04cbb97973ab2e2b9936d8425d"),
			instanceData(8, "af2aaccb6fcdb47197c1c8d1ce202ae44fa6787bae12f55d677e7722152ac193"),
			instanceData(9, "ab9ce6f1aa909f6c1efaa3859604b13aea6b699674226363fd99767d39518bed"),
			instanceData(10, "7cd71f9ecb9e2d63f1128a16a2c4a65d9c4979717b63cd4021a0fbde6476e810"),
		},
	}
}
