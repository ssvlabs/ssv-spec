package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerEnvelopeOnlyPacket tests that the block root must be present in a post-consensus packet (SIP #94
// §4): the §6 envelope entry is the optional one. A peer's packet carrying the envelope entry alone is rejected
// whole, so its envelope share is not collected either.
func GloasProposerEnvelopeOnlyPacket() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer envelope only packet",
		Documentation: testdoc.GloasProposerEnvelopeOnlyPacketDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages: append(
			testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // pre-consensus + consensus
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerEnvelopeOnlyMsgV(ks.Shares[2], 2, version))),
		),
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
		},
		ExpectedErrorCode: types.WrongRootsCountErrorCode,
	}
}
