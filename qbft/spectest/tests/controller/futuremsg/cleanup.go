package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// Cleanup tests cleaning up future msgs container
func Cleanup() *ControllerSyncSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	signMsgEncoded, _ := testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
		Height: 5,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
		Height: 10,
		Round:  3,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
		Height: 11,
		Round:  3,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()

	multiSignMsgEncoded, _ := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: 10,
			Round:  qbft.FirstRound,
			Input:  inputData,
		}).Encode()

	return &ControllerSyncSpecTest{
		Name: "future msgs cleanup",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
				Data: signMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
				Data: signMsgEncoded2,
			},
			{
				ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
				Data: multiSignMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
				Data: signMsgEncoded3,
			},
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "44f061128fec7e83fcdd8cc82300306ddf6ee4f71f4802b37993710c761ae3ba",
	}
}
