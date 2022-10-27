package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastRoundProposalPrevNotPrepared tests a valid proposal for past round (not prev prepared)
func PastRoundProposalPrevNotPrepared() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 10

	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	})
	signMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})

	rcMsgHeader, _ := rcMsg.ToSignedMessage()
	rcMsgHeader2, _ := rcMsg2.ToSignedMessage()
	rcMsgHeader3, _ := rcMsg3.ToSignedMessage()

	signMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsgHeader,
		rcMsgHeader2,
		rcMsgHeader3,
	}
	proposeMsgEncoded, _ := signMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "proposal past round (not prev prepared)",
		Pre:            pre,
		PostRoot:       "36ff2e09d3d5f31610fefda7e7193fc190ebea5bd058b5a163f0c9fd7e5d78e0",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
		ExpectedError:  "proposal invalid: proposal is not valid with current state",
	}
}
