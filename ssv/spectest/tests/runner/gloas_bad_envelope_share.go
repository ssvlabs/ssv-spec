package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerBadEnvelopeShare tests that a bad §6 envelope share does not strand the block (SIP #94 §4/§6).
// Operator 3's packet crosses both roots to quorum with an invalid envelope share, listing the envelope entry
// first; the block still submits, while the envelope fails reconstruction and is not published.
func GloasProposerBadEnvelopeShare() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	badPacket := testingutils.PostConsensusProposerBadEnvelopeShareMsgV(ks.Shares[3], 3, version)
	badPacket.Messages[0], badPacket.Messages[1] = badPacket.Messages[1], badPacket.Messages[0] // envelope entry first

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer bad envelope share",
		Documentation: testdoc.GloasProposerBadEnvelopeShareDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages: append(
			testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // pre-consensus + consensus
			[]*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, badPacket)),
			}...,
		),
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
		},
		ExpectedErrorCode: types.ReconstructSignatureErrorCode,
	}
}
