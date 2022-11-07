package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// GetRoot tests GetRoot on SignedMessage
func GetRoot() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	proposalMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, inputData)

	j := proposalMsg.ToJustification()
	proposalMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		j,
	}
	proposalMsg.ProposalJustifications = []*qbft.SignedMessage{
		j,
	}
	proposalMsgEncoded, _ := proposalMsg.Encode()

	r, _ := proposalMsg.GetRoot()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
			Data: proposalMsgEncoded,
		},
	}

	return &tests.MsgSpecTest{
		Name:     "get root",
		Messages: msgs,
		ExpectedRoots: [][]byte{
			r,
		},
	}
}
