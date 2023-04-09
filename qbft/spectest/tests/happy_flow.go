package tests

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func happyFlowPostState() *qbft.State {
	return &qbft.State{
		Share: &types.Share{
			OperatorID: 1111,
		},
	}
}

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),

		testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingCommitMessage(ks.Shares[3], types.OperatorID(3)),
	}

	post := happyFlowPostState()
	r, err := post.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	RootRegister[hex.EncodeToString(r)] = post

	return &MsgProcessingSpecTest{
		Name:          "happy flow",
		Pre:           pre,
		PostRoot:      hex.EncodeToString(r), //"a298f278a78362257e233e1db8eeb44c2bb9b45a55bb3555928d6723231ebcd2",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
		},
	}
}
