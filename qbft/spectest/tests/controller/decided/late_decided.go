package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateDecided tests processing a decided msg for a just decided instance
func LateDecided() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	msgs := testingutils.DecidingMsgsForHeight(inputData, identifier, qbft.FirstHeight, ks)
	multiSignMsg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
			Input:  inputData,
		})
	multiSignMsg2 := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[4]},
		[]types.OperatorID{1, 2, 4},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
			Input:  inputData,
		})
	multiSignMsgEncoded2, _ := multiSignMsg2.Encode()

	msgs = append(msgs, &types.Message{
		ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
		Data: multiSignMsgEncoded2,
	})
	return &tests.ControllerSpecTest{
		Name: "decide late decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         inputData,
				InputMessages:      msgs,
				SavedDecided:       multiSignMsg,
				BroadcastedDecided: multiSignMsg,
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "d5d4696d29f1359a0f55292ba42dfd922993408529aa86926243df2221554c11",
			},
		},
	}
}
