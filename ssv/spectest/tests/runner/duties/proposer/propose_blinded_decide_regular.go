package proposer

import (
	"encoding/hex"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func getSSZRootNoError(obj ssz.HashRoot) string {
	r, _ := obj.HashTreeRoot()
	return hex.EncodeToString(r[:])
}

// ProposeBlindedBlockDecidedRegular tests proposing a blinded block but the decided block is a regular block. Full flow
func ProposeBlindedBlockDecidedRegular() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MsgProcessingSpecTest{
		Name:   "propose blinded decide regular",
		Runner: testingutils.ProposerBlindedBlockRunner(ks),
		Duty:   testingutils.TestingProposerDuty,
		Messages: []*types.SSVMessage{
			testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
			testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

			testingutils.SSVMsgProposer(
				testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: testingutils.ProposerMsgID,
						Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
					}), nil),

			testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
			testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
			testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
		},
		PostDutyRunnerStateRoot: "cb96348639218e73f803afa8a243d1ccb97a1190fa8fec83d2dbac14cd05849d",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
			testingutils.PostConsensusProposerMsg(testingutils.Testing4SharesSet().Shares[1], 1),
		},
		BeaconBroadcastedRoots: []string{
			getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
		},
	}
}
