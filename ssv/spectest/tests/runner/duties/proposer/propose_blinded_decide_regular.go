package proposer

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func getSSZRootNoError(obj ssz.HashRoot) string {
	if obj == nil {
		return ""
	}
	r, _ := obj.HashTreeRoot()
	return hex.EncodeToString(r[:])
}

// ProposeBlindedBlockDecidedRegular tests proposing a blinded block but the decided block is a regular block. Full flow
func ProposeBlindedBlockDecidedRegular() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "propose blinded decide regular",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "bellatrix",
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
							testingutils.TestProposerBlindedBlockConsensusDataByts(spec.DataVersionBellatrix),
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.ProposerMsgID,
								Root:       sha256.Sum256(testingutils.TestProposerBlindedBlockConsensusDataByts(spec.DataVersionBellatrix)),
							}), nil),

					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1, spec.DataVersionBellatrix)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2, spec.DataVersionBellatrix)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3, spec.DataVersionBellatrix)),
				},
				PostDutyRunnerStateRoot: "a8c49150bc532658895b816aef239a1f2d4a6d97dac9424407ecf60a5dee1676",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
					testingutils.PostConsensusProposerMsg(testingutils.Testing4SharesSet().Shares[1], 1, spec.DataVersionBellatrix),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlock(ks, spec.DataVersionBellatrix)),
				},
			},
			/*{
				Name:   "capella",
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
							testingutils.TestProposerBlindedBlockConsensusDataByts(spec.DataVersionCapella),
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.ProposerMsgID,
								Root:       sha256.Sum256(testingutils.TestProposerBlindedBlockConsensusDataByts(spec.DataVersionCapella)),
							}), nil),

					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1, spec.DataVersionCapella)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2, spec.DataVersionCapella)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3, spec.DataVersionCapella)),
				},
				PostDutyRunnerStateRoot: "cb96348639218e73f803afa8a243d1ccb97a1190fa8fec83d2dbac14cd05849d",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
					testingutils.PostConsensusProposerMsg(testingutils.Testing4SharesSet().Shares[1], 1, spec.DataVersionCapella),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlock(ks, spec.DataVersionCapella)),
				},
				ExpectedError: "failed processing post consensus message: invalid post-consensus message: no decided value", //TODO no ssz hash for blinded block. need to update test to not expect error
			},*/
			{
				Name:   "unknown",
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
							testingutils.TestProposerBlindedBlockConsensusDataByts(spec.DataVersionPhase0),
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.ProposerMsgID,
								Root:       sha256.Sum256(testingutils.TestProposerBlindedBlockConsensusDataByts(spec.DataVersionPhase0)), // mock not supported version
							}), nil),

					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1, spec.DataVersionPhase0)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2, spec.DataVersionPhase0)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3, spec.DataVersionPhase0)),
				},
				PostDutyRunnerStateRoot: "e6226859d56149d07116b587b08dde04157b3c32d82492ed792189596d7c0ac9",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{},
				ExpectedError:          "failed processing post consensus message: invalid post-consensus message: no decided value",
			},
		},
	}
}
