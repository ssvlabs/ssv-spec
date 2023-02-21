package proposer

import (
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// ProposeRegularBlockDecidedBlinded tests proposing a regular block but the decided block is a blinded block. Full flow
func ProposeRegularBlockDecidedBlinded() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "propose regular decide blinded",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "bellatrix",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

					testingutils.SSVMsgProposer(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							testingutils.TestProposerConsensusDataByts(spec.DataVersionBellatrix),
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.ProposerMsgID,
								Root:       sha256.Sum256(testingutils.TestProposerConsensusDataByts(spec.DataVersionBellatrix)),
							}), nil),

					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1, spec.DataVersionBellatrix)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2, spec.DataVersionBellatrix)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3, spec.DataVersionBellatrix)),
				},
				PostDutyRunnerStateRoot: "41cf66d9f00c96d6f9d4d881437a563bcda9c7e7854dc5d6dcef920705a48841",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
					testingutils.PostConsensusProposerMsg(testingutils.Testing4SharesSet().Shares[1], 1, spec.DataVersionBellatrix),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks, spec.DataVersionBellatrix)),
				},
			},
			{
				Name:   "capella",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

					testingutils.SSVMsgProposer(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							testingutils.TestProposerConsensusDataByts(spec.DataVersionCapella),
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.ProposerMsgID,
								Root:       sha256.Sum256(testingutils.TestProposerConsensusDataByts(spec.DataVersionCapella)),
							}), nil),

					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1, spec.DataVersionCapella)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2, spec.DataVersionCapella)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3, spec.DataVersionCapella)),
				},
				PostDutyRunnerStateRoot: "c86a80dfac04c9151f740b38a9191666e5f17186b3b720741a236d1f2b0f94c1",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
					testingutils.PostConsensusProposerMsg(testingutils.Testing4SharesSet().Shares[1], 1, spec.DataVersionCapella),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks, spec.DataVersionCapella)),
				},
			},
			{
				Name:   "unknown",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

					testingutils.SSVMsgProposer(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							testingutils.TestProposerConsensusDataByts(spec.DataVersionPhase0),
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.ProposerMsgID,
								Root:       sha256.Sum256(testingutils.TestProposerConsensusDataByts(spec.DataVersionPhase0)), // mock not supported version
							}), nil),

					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1, spec.DataVersionPhase0)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2, spec.DataVersionPhase0)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3, spec.DataVersionPhase0)),
				},
				PostDutyRunnerStateRoot: "b9d4e4a71ab59c4714a924349ec8de3ad308924b0c00845601a5402f7c6a731b",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{},
				ExpectedError:          "failed processing post consensus message: invalid post-consensus message: no decided value",
			},
		},
	}
}
