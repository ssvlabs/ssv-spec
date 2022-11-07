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
	proposalMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}})

	j := proposalMsg.ToJustification()
	proposalMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		j,
	}
	proposalMsg.ProposalJustifications = []*qbft.SignedMessage{
		j,
	}
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
