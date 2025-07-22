package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnmarshalJustifications tests unmarshalling justifications
func UnmarshalJustifications() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndFullData(ks.OperatorKeys[1], 1, 2, nil),
	}

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithRoundAndFullData(ks.OperatorKeys[1], types.OperatorID(1), 1, nil),
	}

	msg := testingutils.TestingProposalMessageWithParams(
		ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs))

	r, err := msg.GetRoot()
	if err != nil {
		panic(err)
	}

	b, err := msg.Encode()
	if err != nil {
		panic(err)
	}

	return tests.NewMsgSpecTest(
		"unmarshal justifications",
		testdoc.MessagesUnmarshalJustificationDoc,
		[]*types.SignedSSVMessage{msg},
		[][]byte{b},
		[][32]byte{r},
		"",
	)
}
