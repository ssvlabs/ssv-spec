package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// Cleanup tests cleaning up future msgs container
func Cleanup() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "future msgs cleanup",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingCommitMessageWithParams(
				ks.Shares[4], 4, qbft.FirstRound, 5, identifier[:], testingutils.TestingQBFTRootData),
			testingutils.TestingPrepareMessageWithParams(
				ks.Shares[3], 3, 3, 10, identifier[:], testingutils.TestingQBFTRootData),
			testingutils.TestingCommitMultiSignerMessageWithParams(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				qbft.FirstRound,
				10,
				identifier[:],
				testingutils.TestingQBFTRootData,
				testingutils.TestingQBFTFullData,
			),
			testingutils.TestingPrepareMessageWithParams(
				ks.Shares[2], 2, 3, 11, identifier[:], testingutils.TestingQBFTRootData),
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "cd0e7827cc4f55c6972925aa38b79050ae7f99d6083032127b9ee1fbd28caae0",
	}
}
