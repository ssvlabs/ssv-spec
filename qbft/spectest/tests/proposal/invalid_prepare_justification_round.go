package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidPrepareJustificationRound tests a proposal for > 1 round, prepared previously but one of the prepare justifications has round != highest prepared round
func InvalidPrepareJustificationRound() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	prepareMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	prepareMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	prepareMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{Root: pre.StartValue.Root})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		PreparedRound: qbft.FirstRound,
	}, pre.StartValue)
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		PreparedRound: qbft.FirstRound,
	}, pre.StartValue)
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		PreparedRound: qbft.FirstRound,
	}, pre.StartValue)
	proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, pre.StartValue)

	justifications := []*qbft.SignedMessage{
		prepareMsg,
		prepareMsg2,
		prepareMsg3,
	}

	rcMsg.RoundChangeJustifications = justifications
	rcMsg2.RoundChangeJustifications = justifications
	rcMsg3.RoundChangeJustifications = justifications

	proposeMsg.ProposalJustifications = justifications
	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg.ToJustification(),
		rcMsg2.ToJustification(),
		rcMsg3.ToJustification(),
	}
	proposeMsgEncoded, _ := proposeMsg.Encode()

	return &tests.MsgProcessingSpecTest{
		Name:     "invalid prepare justification round",
		Pre:      pre,
		PostRoot: "56cee2fd474513bc56851dfbb027366f6fc3f90fe8fec4081e993b69f84e2228",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposeMsgEncoded,
			},
		},
		OutputMessages: []*types.Message{},
		ExpectedError:  "proposal invalid: proposal not justified: change round msg not valid: round change justification invalid: msg round wrong",
	}
}
