package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateSigner tests multiple future msg for the same signer (doesn't trigger futuremsg)
func DuplicateSigner() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()

	return &ControllerSyncSpecTest{
		Name: "future msg duplicate signer",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingCommitMessageWithHeight(ks.Shares[4], 4, 5),
			testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 3, 3, 10, testingutils.TestingQBFTRootData),

			testingutils.TestingPrepareMessageWithHeight(ks.Shares[4], 4, 6),
			testingutils.TestingRoundChangeMessageWithHeight(ks.Shares[4], 4, 2),
			testingutils.TestingCommitMessageWithHeight(ks.Shares[4], 4, 50),
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "ca7eaf0f0b404b601dc8dd471924794ce32ef6bcb88721098b7b6014001754c1",
		ExpectedError:        "discarded future msg",
	}
}
