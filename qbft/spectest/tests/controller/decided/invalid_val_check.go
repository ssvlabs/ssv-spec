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
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{
		Root:   [32]byte{1, 2, 3, 4},
		Source: []byte{1, 2, 3, 4},
	}
	multiMsg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: 10,
			Round:  qbft.FirstRound,
		}, &qbft.Data{
			Root:   [32]byte{1, 1, 1, 1},
			Source: testingutils.TestingInvalidValueCheck,
		})
	multiMsgEncoded, _ := multiMsg.Encode()
	return &tests.ControllerSpecTest{
		Name: "decide invalid value (should pass)",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					{
						ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
						Data: multiMsgEncoded,
					},
				},
				SavedDecided:       multiMsg,
				DecidedVal:         testingutils.TestingInvalidValueCheck,
				DecidedCnt:         1,
				ControllerPostRoot: "09c3550c6fcbefb63ad2a66013c693fe0f4e602d5f05214532ab1acad8213e65",
			},
		},
	}
}
