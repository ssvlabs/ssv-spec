package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
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

	test := tests.NewMsgProcessingSpecTest(
		"duplicate rc msg justification",
		testdoc.ProposalDuplicateRCMsgDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		types.RoundChangeNoQuorumErrorCode,
		nil,
		ks,
	)

	return test
}
