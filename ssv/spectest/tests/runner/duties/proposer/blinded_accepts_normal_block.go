package proposer

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BlindedRunnerAcceptsNormalBlock tests a compatibility flow where a proposer
// runner that now starts QBFT with blinded data still accepts a full-block
// decision from the network. This models mixed-version clusters during the
// transition from full proposer values to blinded-only proposer values.
func BlindedRunnerAcceptsNormalBlock() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"blinded proposer accepts normal block proposal",
		testdoc.ProposerBlindedReceivingNormalBlockDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	// proposerBlindedReceivingNormalBlockV creates a compatibility test
	// specification for a versioned blinded proposer receiving a full-block
	// consensus decision from peers that still originate the old value form.
	proposerBlindedReceivingNormalBlockV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("blinded proposer accepts normal block proposal (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // consensus
				[]*types.SignedSSVMessage{ // post consensus
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
				}...,
			),
			PostDutyRunnerStateRoot: fullHappyFlowBlindedProposerReceivingNormalBlockSC(version).Root(),
			PostDutyRunnerState:     fullHappyFlowBlindedProposerReceivingNormalBlockSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerBlindedReceivingNormalBlockV(v)}...)
	}

	return multiSpecTest
}
