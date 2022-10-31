package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoPrevAcceptedProposal tests a commit msg received without a previous accepted proposal
func NoPrevAcceptedProposal() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = nil
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: pre.StartValue.Root},
	}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "no previous accepted proposal",
		Pre:           pre,
		PostRoot:      "56cee2fd474513bc56851dfbb027366f6fc3f90fe8fec4081e993b69f84e2228",
		InputMessages: msgs,
		ExpectedError: "did not receive proposal for this round",
	}
}
