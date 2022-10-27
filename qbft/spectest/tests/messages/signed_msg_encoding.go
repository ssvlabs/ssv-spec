package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignedMessageEncoding tests encoding SignedMessage
func SignedMessageEncoding() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	proposalMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})

	prepareMsgHeader, _ := signQBFTMsg.ToSignedMessage()
	prepareJustifications := []*qbft.SignedMessage{
		prepareMsgHeader,
	}
	proposalMsg.RoundChangeJustifications = prepareJustifications
	proposalMsg.ProposalJustifications = prepareJustifications

	proposalMsgEncoded, _ := proposalMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
			Data: proposalMsgEncoded,
		},
	}

	return &tests.MsgSpecTest{
		Name:     "signed message encoding",
		Messages: msgs,
		EncodedMessages: [][]byte{
			proposalMsgEncoded,
		},
	}
}
