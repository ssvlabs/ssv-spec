package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongData tests prepare msg with data != acceptedProposalData.Data
func WrongData() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	msgs := []*qbft.SignedMessage{
		// TODO: different value instead of wrong root
		testingutils.TestingPrepareMessageWrongRoot(ks.Shares[1], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare wrong data",
		Pre:           pre,
		PostRoot:      "b61f5233721865ca43afc68f4ad5045eeb123f6e8f095ce76ecf956dabc74713",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposed data mistmatch",
	}
}
