package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateSigner tests multiple future msg for the same signer (doesn't trigger futuremsg)
func DuplicateSigner() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "future msg duplicate signer",
		InputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     5,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Root:       testingutils.TestingQBFTRootData,
			}),
			testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     10,
				Round:      3,
				Identifier: identifier[:],
				Root:       testingutils.TestingQBFTRootData,
			}),

			testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     6,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Root:       testingutils.TestingQBFTRootData,
			}),
			testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
				MsgType:    qbft.RoundChangeMsgType,
				Height:     2,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Root:       testingutils.TestingQBFTRootData,
			}),
			testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     50,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Root:       testingutils.TestingQBFTRootData,
			}),
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "98c0625374dc64e6350eb704812d9222f9a6121a87c2a55a8b1a3f8790e87c77",
		ExpectedError:        "discarded future msg",
	}
}
