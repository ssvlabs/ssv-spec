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
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 10),
	}

	test := tests.NewCreateMsgSpecTest(
		"create proposal not previously prepared",
		testdoc.MessagesCreateProposalNotPrevPreparedDoc,
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData,
		10,
		roundChangeJustifications,
		nil,
		tests.CreateProposal,
		"ebe29a35a3862c7f720568f6aea8273e522d2a4307e84eb08b91fe2fbd8a2920",
		nil,
		"",
		ks,
	)

	return test
}
