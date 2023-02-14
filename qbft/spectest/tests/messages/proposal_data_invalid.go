package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ProposalDataInvalid tests proposal data is invalid
func ProposalDataInvalid() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), qbft.FirstRound, qbft.FirstHeight, testingutils.WrongRoot,
		nil, nil)

	return &tests.MsgSpecTest{
		Name: "proposal data invalid",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "ProposalData data is invalid",
	}
}
