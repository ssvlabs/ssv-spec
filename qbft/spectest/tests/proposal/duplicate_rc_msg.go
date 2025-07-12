package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicateRCMsg tests a proposal for > 1 round, not prepared previously with quorum of round change but 2 are duplicates (shouldn't find quorum)
func DuplicateRCMsg() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
	}

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, testingutils.MarshalJustifications(rcMsgs), nil),
	}

	return tests.NewMsgProcessingSpecTest(
		"duplicate rc msg justification",
		"Test proposal for round > 1 that was not prepared previously with round change messages but contains duplicates, expecting validation error due to no quorum.",
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: proposal not justified: change round has no quorum",
		nil,
	)
}
