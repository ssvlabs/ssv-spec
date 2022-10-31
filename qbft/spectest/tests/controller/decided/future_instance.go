package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// FutureInstance tests decided msg received for future (not yet started) instance
func FutureInstance() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	multiSignMsg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: 10,
			Round:  qbft.FirstRound,
			Input:  inputData,
		})
	multiSignMsgEncoded, _ := multiSignMsg.Encode()
	return &tests.ControllerSpecTest{
		Name: "decide future instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					{
						ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
						Data: multiSignMsgEncoded,
					},
				},
				SavedDecided:       multiSignMsg,
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "cf79aafb1c5b30c9c737026d3263934ee31da72db9b68bedc21c5e0f7def69b0",
			},
		},
	}
}
