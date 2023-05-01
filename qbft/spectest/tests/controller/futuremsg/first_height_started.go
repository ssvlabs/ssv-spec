package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FirstHeightStarted tests a future message for special case (first height, instance started)
func FirstHeightStarted() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "future message first height started",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingProposalMessageWithID(ks.Shares[1], 1, identifier),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "783e12e3ddf0f09930b4ce1d1feae7e5ca29133037c81fa3b0ee534ce52294a0",
	}
}
