package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// DuplicateMsgDifferentRoot tests a duplicate proposal msg processing (second one with different root)
func DuplicateMsgDifferentRoot() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := duplicateMsgDifferentRootStateComparison()
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingProposalMessageDifferentRoot(ks.OperatorKeys[1], types.OperatorID(1)),
	}
	outputMessages := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	test := tests.NewMsgProcessingSpecTest(
		"proposal duplicate message different value",
		testdoc.ProposalDuplicateMsgDifferentRootDoc,
		pre,
		sc.Root(),
		sc.ExpectedState,
		inputMessages,
		outputMessages,
		"invalid signed message: proposal is not valid with current state",
		nil,
	)

	test.SetPrivateKeys(ks)

	return test
}

func duplicateMsgDifferentRootStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	instance := &qbft.Instance{
		State: &qbft.State{
			CommitteeMember:                 testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
			ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))),
			ID:                              testingutils.TestingIdentifier,
			Round:                           qbft.FirstRound,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	return &comparable.StateComparison{ExpectedState: instance.State}
}
