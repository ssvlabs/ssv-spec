package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ProposeDataEncoding tests encoding ProposalData
func ProposeDataEncoding() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}})

	j := msg.ToJustification()
	msg.RoundChangeJustifications = []*qbft.SignedMessage{
		j,
	}
	msg.ProposalJustifications = []*qbft.SignedMessage{
		j,
	}

	r, _ := msg.GetRoot()
	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "propose data encoding",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
				Data: b,
			},
		},
		EncodedMessages: [][]byte{
			b,
		},
		ExpectedRoots: [][]byte{
			r,
		},
	}
}
