package proposer

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BlindedRunnerAcceptsNormalBlock tests a full happy flow for a proposer runner that produces blinded blocks but receives as proposal a normal block
func BlindedRunnerAcceptsNormalBlock() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name:  "blinded proposer accepts normal block proposal",
		Tests: []*tests.MsgProcessingSpecTest{},
	}

	// proposerBlindedReceivingNormalBlockV creates a test specification for versioned blinded proposer receiving a normal block.
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
