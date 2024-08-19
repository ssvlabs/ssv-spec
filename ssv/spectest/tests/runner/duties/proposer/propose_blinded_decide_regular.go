package proposer

import (
	"crypto/rsa"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/qbft"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ProposeBlindedBlockDecidedRegular tests proposing a blinded block but the decided block is a regular block. Full flow
func ProposeBlindedBlockDecidedRegular() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MsgProcessingSpecTest{
		Name:   "propose blinded decide regular",
		Runner: testingutils.ProposerBlindedBlockRunner(ks),
		Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionDeneb))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, spec.DataVersionDeneb))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, spec.DataVersionDeneb))),

			testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
				[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
				[]types.OperatorID{1, 2, 3},
				qbft.Height(testingutils.TestingDutySlotV(spec.DataVersionDeneb)),
				testingutils.ProposerMsgID,
				testingutils.TestProposerConsensusDataBytsV(spec.DataVersionDeneb),
			),

			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb))),
		},
		PostDutyRunnerStateRoot: "05c3df4f48431ba9cf2b410358300a01aaae16176f73a4ba192a9d8ce327fba9",
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
			testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, spec.DataVersionDeneb)),
		},
	}
}
