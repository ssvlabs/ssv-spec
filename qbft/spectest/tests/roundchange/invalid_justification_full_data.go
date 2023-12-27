package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidJustificationFullData tests round change justification for which H(full data) != root
func InvalidJustificationFullData() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	invalidFullData := make([]byte, 40)

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 3),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 3),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 3),
	}
	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundRCDataRoundAndFullData(ks.Shares[1], types.OperatorID(1), 3,
			testingutils.MarshalJustifications(prepareMsgs), 3, invalidFullData),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "justification invalid full data",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: H(data) != root",
	}
}
