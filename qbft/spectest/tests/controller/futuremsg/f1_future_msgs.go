package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1FutureMsgs tests a f+1 future msgs that trigger decided sync
func F1FutureMsgs() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "f+1 future msgs",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingCommitMessageWithParams(
				ks.Shares[4], 4, qbft.FirstRound, 5, identifier[:], testingutils.TestingQBFTRootData),
			testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 3, 3, 10, identifier[:], testingutils.TestingQBFTRootData),
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "98c0625374dc64e6350eb704812d9222f9a6121a87c2a55a8b1a3f8790e87c77",
	}
}
