package proposer

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// ProposeRegularBlockDecidedBlinded tests proposing a regular block but the decided block is a blinded block. Full flow
func ProposeRegularBlockDecidedBlinded() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MsgProcessingSpecTest{
		Name:   "propose regular decide blinded",
		Runner: testingutils.ProposerRunner(ks),
		Duty:   &testingutils.TestingProposerDuty,
		Messages: []*types.SSVMessage{
			testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
			testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

			testingutils.SSVMsgProposer(
				testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					testingutils.ProposerMsgID,
					testingutils.TestProposerBlindedBlockConsensusDataByts,
				), nil),

			testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
			testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
			testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
		},
		PostDutyRunnerStateRoot: "d3a677e479864a7ec27e9315cc62fde5eb7f4081b47461b4670a497a91ecce7f",
		OutputMessages: []*types.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
			testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
		},
		BeaconBroadcastedRoots: []string{
			getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
		},
	}
}
