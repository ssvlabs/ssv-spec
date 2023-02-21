package dutyexe

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongDutyPubKey tests decided value with duty validator pubkey != the duty runner's pubkey
func WrongDutyPubKey() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAttester(testingutils.TestingProposalMessage(ks.Shares[1], 1), nil),

		testingutils.SSVMsgAttester(testingutils.TestingPrepareMessage(ks.Shares[1], 1), nil),
		testingutils.SSVMsgAttester(testingutils.TestingPrepareMessage(ks.Shares[2], 2), nil),
		testingutils.SSVMsgAttester(testingutils.TestingPrepareMessage(ks.Shares[3], 3), nil),

		testingutils.SSVMsgAttester(testingutils.TestingCommitMessage(ks.Shares[1], 1), nil),
		testingutils.SSVMsgAttester(testingutils.TestingCommitMessage(ks.Shares[2], 2), nil),
		testingutils.SSVMsgAttester(testingutils.TestingCommitMessage(ks.Shares[3], 3), nil),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "wrong decided value's pubkey",
		Runner:                  dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "3f82ba9763ce97791e62f6daff599692f82608dbff222e8f6562a48a34f08272",
		ExpectedError:           "decided value is invalid: decided value's validator pk is wrong",
	}
}
