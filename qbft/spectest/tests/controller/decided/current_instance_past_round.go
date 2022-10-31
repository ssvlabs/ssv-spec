package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// CurrentInstancePastRound tests a decided msg received for current running instance for a past round
func CurrentInstancePastRound() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  inputData,
	}).Encode()
	proposeMsg2 := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  inputData,
	})
	signedMsgEncoded, _ := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	signedMsgEncoded2, _ := testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	signedMsgEncoded4, _ := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	signedMsgEncoded5, _ := testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	signedMsgEncoded6, _ := testingutils.SignQBFTMsg(ks.Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	rcMsg := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{},
	})
	rcMsg2 := testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{},
	})
	rcMsg3 := testingutils.SignQBFTMsg(ks.Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{},
	})
	proposeMsg2.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg, rcMsg2, rcMsg3,
	}
	rcMsgEncoded, _ := rcMsg.Encode()
	rcMsgEncoded2, _ := rcMsg2.Encode()
	rcMsgEncoded3, _ := rcMsg3.Encode()
	proposeMsgEncoded2, _ := proposeMsg2.Encode()
	multiSignMsg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
			Input:  inputData,
		})
	multiSignMsgEncoded, _ := multiSignMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
			Data: signedMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
			Data: signedMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded3,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
			Data: signedMsgEncoded4,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
			Data: signedMsgEncoded5,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
			Data: signedMsgEncoded6,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
			Data: signedMsgEncoded5,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
			Data: signedMsgEncoded6,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
			Data: multiSignMsgEncoded,
		},
	}

	return &tests.ControllerSpecTest{
		Name: "decide current instance past round",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         inputData,
				InputMessages:      msgs,
				SavedDecided:       multiSignMsg,
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "a8dedba7fa4094b697968c9fd6ba08ec4aaadce5d79077201b72595bc586389f",
			},
		},
	}
}
