package proposer

import (
	"crypto/rsa"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ProposeRegularBlockDecidedBlinded tests proposing a regular block but the decided block is a blinded block. Full flow
func ProposeRegularBlockDecidedBlinded() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	test := tests.NewMsgProcessingSpecTest(
		"propose regular decide blinded",
		testdoc.ProposeRegularBlockDecidedBlindedDoc,
		testingutils.ProposerRunner(ks),
		testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
		[]*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionDeneb))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, spec.DataVersionDeneb))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, spec.DataVersionDeneb))),

			testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
				[]*rsa.PrivateKey{
					ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3],
				},
				[]types.OperatorID{1, 2, 3},
				qbft.Height(testingutils.TestingDutySlotV(spec.DataVersionDeneb)),
				testingutils.ProposerMsgID,
				testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
			),

			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb))),
		},
		false,
		"97b5fd3d786658f67d9d39df63a7a73b690e7873f9bdc107f6fcd401a42d98fc",
		nil,
		[]*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
			testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
		},
		[]string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlockV(ks, spec.DataVersionDeneb)),
		},
		false,
		0,
		ks,
	)

	return test
}
