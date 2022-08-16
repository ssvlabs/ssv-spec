package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastRound tests a proposal for signedProposal.Message.Round < state.Round
func PastRound() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 5
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      5,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
	})
	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data: testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, []*qbft.SignedMessage{
				testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
					MsgType:    qbft.RoundChangeMsgType,
					Height:     qbft.FirstHeight,
					Round:      2,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
				}),
				testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
					MsgType:    qbft.RoundChangeMsgType,
					Height:     qbft.FirstHeight,
					Round:      2,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
				}),
				testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
					MsgType:    qbft.RoundChangeMsgType,
					Height:     qbft.FirstHeight,
					Round:      2,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
				}),
			}, nil),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "proposal past round",
		Pre:            pre,
		PostRoot:       "02c53d76bdfa84c573386a7dff3e443f120d441b3086f7d5e3834a5c7e1261ab",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "proposal invalid: proposal is not valid with current state",
	}
}
