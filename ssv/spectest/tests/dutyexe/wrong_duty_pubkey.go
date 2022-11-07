package dutyexe

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongDutyPubKey tests decided value with duty validator pubkey != the duty runner's pubkey
func WrongDutyPubKey() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	dr := testingutils.AttesterRunner(ks)
	inputData := &qbft.Data{Root: testingutils.TestConsensusWrongDutyPKDataRoot, Source: testingutils.TestConsensusWrongDutyPKDataByts}

	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, inputData).Encode()
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
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
			ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
			Data: signMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
			Data: signMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
			Data: signMsgEncoded3,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "wrong decided value's pubkey",
		Runner:                  dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "3f82ba9763ce97791e62f6daff599692f82608dbff222e8f6562a48a34f08272",
		ExpectedError:           "decided value is invalid: decided value's validator pk is wrong",
	}
}
