package roundchange

import (
	"crypto/sha256"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// QuorumNotTriggeredTwiceJustificationIgnored tests that the fourth round change message does not trigger a quorum and
// no proposal is sent. Also, the justification is ignored.
func QuorumNotTriggeredTwiceJustificationIgnored() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = nil // proposal resets on upon timeout
	pre.State.Round = 2
	testData := []byte{1, 2}
	testDataRoot := sha256.Sum256(testData)

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithFullData(ks.OperatorKeys[1], types.OperatorID(1), testData),
		testingutils.TestingPrepareMessageWithFullData(ks.OperatorKeys[2], types.OperatorID(2), testData),
		testingutils.TestingPrepareMessageWithFullData(ks.OperatorKeys[3], types.OperatorID(3), testData),
	}
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
		testingutils.TestingRoundChangeMessageWithParamsAndFullData(ks.OperatorKeys[4], types.OperatorID(4), 2, qbft.FirstHeight,
			testDataRoot, 1, testingutils.MarshalJustifications(prepareMsgs), testData),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "quorum not triggered twice justification ignored",
		Pre:           pre,
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
				testingutils.TestingQBFTRootData, testingutils.MarshalJustifications(msgs[:len(msgs)-1]),
				[][]byte{}),
		},
	}
}
