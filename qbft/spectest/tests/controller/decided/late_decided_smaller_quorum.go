package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateDecidedSmallerQuorum tests processing a decided msg for a just decided instance (with a smaller quorum)
func LateDecidedSmallerQuorum() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
	multiSignMsg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3], ks.Shares[4]},
		[]types.OperatorID{1, 2, 3, 4},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
		}, inputData)
	multiSignMsg2 := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
		}, inputData)
	multiSignMsgEncoded, _ := multiSignMsg.Encode()
	multiSignMsgEncoded2, _ := multiSignMsg2.Encode()
	return &tests.ControllerSpecTest{
		Name: "decide late decided smaller quorum",
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
				ControllerPostRoot: "70c5218e3832249ba51590e5e8850d60863e3a2281e090669448075e68795a25",
			},
		},
	}
}
