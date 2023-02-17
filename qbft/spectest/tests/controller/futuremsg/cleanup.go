package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// Cleanup tests cleaning up future msgs container
func Cleanup() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "future msgs cleanup",
		InputMessages: []*qbft.SignedMessage{
			// TODO: create helper functions receiving identifiers for all futuremsg tests
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
			testingutils.TestingCommitMultiSignerMessageWithParams(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				qbft.FirstRound,
				10,
				identifier[:],
				testingutils.TestingQBFTRootData,
				testingutils.TestingQBFTFullData,
			),
			testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     11,
				Round:      3,
				Identifier: identifier[:],
				Root:       testingutils.TestingQBFTRootData,
			}),
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "2a506a58b7abbbcffd6d4ee92f154c826947ce56672564155772435456315c6d",
	}
}
