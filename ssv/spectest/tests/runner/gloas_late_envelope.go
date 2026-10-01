package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerLateEnvelopeQuorum tests that a Gloas (ePBS) self-build proposer publishes the §6 reveal
// when the envelope root reaches quorum only after the block root (SIP #94 §4). The block root is the
// required entry, so its quorum finishes the duty; the optional envelope root reconstructs independently
// and may complete later. Here operator 3 signs block-only, so the block finishes on its packet (block
// shares 1–3) while the envelope holds two shares (1, 2); operator 4's packet then completes the envelope
// quorum (1, 2, 4). The builder operator — this operator, which produced the decided block — keeps
// accepting post-consensus packets past the block's finish and publishes the reveal on that later quorum,
// so both the block and the envelope are broadcast. Without that allowance operator 4 would be rejected as
// "no running duty" and the reveal missed.
func GloasProposerLateEnvelopeQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer late envelope quorum",
		Documentation: testdoc.GloasProposerLateEnvelopeQuorumDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages: append(
			testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // pre-consensus + consensus
			[]*types.SignedSSVMessage{ // post-consensus: op3 signs block-only, so the envelope quorum lands one packet after the block's
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerBlockOnlyMsgV(ks.Shares[3], 3, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, version))),
			}...,
		),
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
		},
		BeaconBroadcastedRoots: testingutils.WithGloasEnvelopeBroadcast([]string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
		}, version),
	}
}
