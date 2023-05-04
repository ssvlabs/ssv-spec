package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FirstHeightNotStarted tests a future message for special case (first height, instance not started)
func FirstHeightNotStarted() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "future message first height not started",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingProposalMessageWithID(ks.Shares[1], 1, identifier),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "f233b1f7376746543acf153474be2ae28f3c21d6aed52c8b920be7af9297fc66",
		SkipInstanceStart:    true,
	}
}
