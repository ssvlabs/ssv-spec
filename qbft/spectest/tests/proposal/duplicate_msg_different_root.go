package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// DuplicateMsgDifferentRoot tests a duplicate proposal msg processing (second one with different root)
func DuplicateMsgDifferentRoot() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := duplicateMsgDifferentRootStateComparison()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingProposalMessageDifferentRoot(ks.Shares[1], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal duplicate message different value",
		Pre:           pre,
		PostRoot:      sc.Root(),
		PostState:     sc.ExpectedState,
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		},
		ExpectedError: "invalid signed message: proposal is not valid with current state",
	}
}

func duplicateMsgDifferentRootStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
	}

	instance := &qbft.Instance{
		State: &qbft.State{
			Share:                           testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
			ID:                              testingutils.TestingIdentifier,
			Round:                           qbft.FirstRound,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	return &comparable.StateComparison{ExpectedState: instance.State}
}
