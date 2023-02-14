package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ImparsableProposalData tests a proposal msg received with imparsable data
func ImparsableProposalData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Root:       [32]byte{1, 2, 3, 4},
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "imparsable proposal data",
		Pre:           pre,
		PostRoot:      "e9a7c1b25a014f281d084637bdc50ec144f1262b5cc87eeb6b4493d27d10a69b",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: could not get proposal data: could not decode proposal data from message: invalid character '\\x01' looking for beginning of value",
	}
}
