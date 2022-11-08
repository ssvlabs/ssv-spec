package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateDecidedBiggerQuorum tests processing a decided msg for a just decided instance (with a bigger quorum)
func LateDecidedBiggerQuorum() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
	msgs := testingutils.DecidingMsgsForHeight(inputData, identifier, qbft.FirstHeight, ks)
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
	msgs = append(msgs, &types.Message{
		ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
		Data: multiSignMsgEncoded,
	})
	return &tests.ControllerSpecTest{
		Name: "decide late decided bigger quorum",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         inputData,
				InputMessages:      msgs,
				SavedDecided:       multiSignMsg2,
				BroadcastedDecided: multiSignMsg2,
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "f8a7bf2ad195b8945ba74a38739a8d37fe0d88c1d96cda2a30f4a7ade2913391",
			},
		},
	}
}
