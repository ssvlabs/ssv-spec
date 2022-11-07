package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateSigner tests multiple future msg for the same signer (doesn't trigger futuremsg)
func DuplicateSigner() *ControllerSyncSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	signMsgEncoded, _ := testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
		Height: 5,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
		Height: 10,
		Round:  3,
	}, &qbft.Data{Root: inputData.Root}).Encode()

	signMsgEncoded3, _ := testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
		Height: 6,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root}).Encode()
	signMsgEncoded4, _ := testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
		Height: 2,
		Round:  qbft.FirstRound,
	}, inputData).Encode()
	signMsgEncoded5, _ := testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
		Height: 50,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root}).Encode()

	return &ControllerSyncSpecTest{
		Name: "future msg duplicate signer",
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
				ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
				Data: signMsgEncoded3,
			},
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusRoundChangeMsgType),
				Data: signMsgEncoded4,
			},
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
				Data: signMsgEncoded5,
			},
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "871c31e7f01443af45e67b7422cc000c45c9423f1138a761e3d1b306a4f4d78a",
		ExpectedError:        "discarded future msg",
	}
}
