package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// FullFlowAfterDecided tests a decided msg for round 1 followed by a full proposal, prepare, commit for round 2
func FullFlowAfterDecided() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	proposeMsg := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  inputData,
	})
	signedMsgEncoded, _ := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	signedMsgEncoded2, _ := testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	signedMsgEncoded3, _ := testingutils.SignQBFTMsg(ks.Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: inputData.Root},
	}).Encode()
	signedMsgEncoded4, _ := testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
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
	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg, rcMsg2, rcMsg3,
	}
	proposeMsgEncoded, _ := proposeMsg.Encode()
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
			ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
			Data: multiSignMsgEncoded,
		},
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
			ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
			Data: signedMsgEncoded3,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
			Data: signedMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
			Data: signedMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
			Data: signedMsgEncoded4,
		},
	}

	return &tests.ControllerSpecTest{
		Name: "full flow after decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         inputData,
				InputMessages:      msgs,
				SavedDecided:       multiSignMsg,
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "15fc26f1ddb1ee3e56b2ea9f27c5d3740ace5c58b1e90340d1d38aafb46b1f58",
			},
		},
	}
}
