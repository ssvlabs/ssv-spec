package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CreateProposalNotPreviouslyPrepared tests creating a proposal msg, non-first round and not previously prepared
func CreateProposalNotPreviouslyPrepared() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	roundChangeJustifications := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
		testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
	}

	test := tests.NewCreateMsgSpecTest(
		"create proposal not previously prepared",
		testdoc.MessagesCreateProposalNotPrevPreparedDoc,
		[32]byte{1, 2, 3, 4},
		nil,
		10,
		roundChangeJustifications,
		nil,
		tests.CreateProposal,
		"6a2917ae827e875a646e88ebb1d483a0a99e4f321e7f063138e99a7e7b08794e",
		nil,
		"",
	)

	test.SetPrivateKeys(ks)

	return test
}
