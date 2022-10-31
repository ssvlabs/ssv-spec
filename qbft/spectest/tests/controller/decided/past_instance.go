package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// PastInstance tests a decided msg received for past instance
func PastInstance() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	multiSignMsg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: 100,
			Round:  qbft.FirstRound,
			Input:  inputData,
		})
	multiSignMsgEncoded, _ := multiSignMsg.Encode()
	multiSignMsg2 := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: 99,
			Round:  qbft.FirstRound,
			Input:  inputData,
		})
	multiSignMsgEncoded2, _ := multiSignMsg2.Encode()
	return &tests.ControllerSpecTest{
		Name: "decide past instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					{
						ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
						Data: multiSignMsgEncoded,
					},
					{
						ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
						Data: multiSignMsgEncoded2,
					},
				},
				SavedDecided:       multiSignMsg,
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "b1446134a653f433b66c37bd20017153e4993c6ab014d2d87c50ed2c0a77e884",
			},
		},
	}
}
