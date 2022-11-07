package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// JustificationsNotHeighest tests a proposal for > 1 round, prepared previously with rc justification prepares at different heights but the prepare justification is not the highest
func JustificationsNotHeighest() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 3

	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	signQBFTMsg4 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{Root: pre.StartValue.Root})
	signQBFTMsg5 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{Root: pre.StartValue.Root})
	signQBFTMsg6 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{Root: pre.StartValue.Root})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         3,
		PreparedRound: qbft.FirstRound,
	}, pre.StartValue)
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         3,
		PreparedRound: 2,
	}, pre.StartValue)
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  3,
	}, &qbft.Data{})
	proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  3,
	}, pre.StartValue)

	rcMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		signQBFTMsg,
		signQBFTMsg2,
		signQBFTMsg3,
	}
	rcMsg2.RoundChangeJustifications = []*qbft.SignedMessage{
		signQBFTMsg4,
		signQBFTMsg5,
		signQBFTMsg6,
	}

	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg,
		rcMsg2,
		rcMsg3,
	}
	proposeMsg.ProposalJustifications = []*qbft.SignedMessage{
		signQBFTMsg,
		signQBFTMsg2,
		signQBFTMsg3,
	}
	proposeMsgEncoded, _ := proposeMsg.Encode()

	return &tests.MsgProcessingSpecTest{
		Name:     "proposal justification not highest",
		Pre:      pre,
		PostRoot: "e9da16ca03e556777e7c02b59c7a490bb3d54e5bd8f3135ed1d1ef7116f85fb4",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposeMsgEncoded,
			},
		},
		OutputMessages: []*types.Message{},
		ExpectedError:  "proposal invalid: proposal not justified: signed prepare not valid",
	}
}
