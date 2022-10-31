package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LatePrepare tests process late prepare msg for an instance which just decided
func LatePrepare() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	msgs := testingutils.DecidingMsgsForHeight(inputData, identifier, qbft.FirstHeight, ks)

	signedMsgEncoded, _ := testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	lateMsg := &types.Message{
		ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
		Data: signedMsgEncoded,
	}
	multiSignMsg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
			Input:  inputData,
		})

	msgs = append(msgs, lateMsg)

	return &tests.ControllerSpecTest{
		Name: "late prepare",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         inputData,
				InputMessages:      msgs,
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				SavedDecided:       multiSignMsg,
				BroadcastedDecided: multiSignMsg,
				ControllerPostRoot: "5b678d1d8273acda5f51f4f4c3c722ef803c66d1ec228cc6ca082fac3a820bb9",
			},
		},
	}
}
