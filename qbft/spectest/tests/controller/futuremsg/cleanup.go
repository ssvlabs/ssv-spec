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

	return &ControllerSyncSpecTest{
		Name: "future msgs cleanup",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingCommitMessageWithHeight(ks.Shares[4], 4, 5),
			testingutils.TestingPrepareMessageWithParams(ks.Shares[4], 4, 3, 10, testingutils.TestingQBFTRootData),
			testingutils.TestingCommitMultiSignerMessageWithHeight(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				10,
			),
			testingutils.TestingPrepareMessageWithParams(ks.Shares[2], 2, 3, 11, testingutils.TestingQBFTRootData),
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "2a506a58b7abbbcffd6d4ee92f154c826947ce56672564155772435456315c6d",
	}
}
