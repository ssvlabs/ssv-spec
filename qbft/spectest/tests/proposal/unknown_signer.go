package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests a single proposal received with an unknown signer
func UnknownSigner() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(5), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, pre.StartValue).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "unknown proposal signer",
		Pre:           pre,
		PostRoot:      "56cee2fd474513bc56851dfbb027366f6fc3f90fe8fec4081e993b69f84e2228",
		InputMessages: msgs,
		ExpectedError: "proposal invalid: proposal msg signature invalid: unknown signer",
	}
}
