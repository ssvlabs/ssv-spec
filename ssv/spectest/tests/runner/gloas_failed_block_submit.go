package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerFailedBlockSubmit tests that the §6 reveal publishes when this operator's own block submit fails:
// it gates on the submit being attempted, since the block still reaches beacon nodes through the other
// operators (SIP #94 §6).
func GloasProposerFailedBlockSubmit() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer failed block submit",
		Documentation: testdoc.GloasProposerFailedBlockSubmitDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages: append(
			testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // pre-consensus + consensus
			[]*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
			}...,
		),
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
		},
		BeaconBroadcastedRoots: testingutils.WithGloasEnvelopeBroadcast(nil, version),
		BeaconNode:             &tests.BeaconNodeBehaviour{FailBlockSubmit: true},
		ExpectedErrorCode:      types.ProposerBlockSubmitFailedErrorCode,
	}
}
