package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerEnvelopeFirstOrder tests that a Gloas (ePBS) self-build proposer submits both the block and
// the §6 envelope when peers' post-consensus packets list the envelope entry before the block entry
// (SIP #94 §4/§6). The runner classifies entries by matched root, not position, and always submits the
// block before publishing the reveal, so packet order does not change the outcome — a guard against a
// future order-dependent rewrite (an early return or positional iteration) reintroducing a strand. The
// full flow runs so the self-build produce sets the envelope the reveal publishes; operators 2 and 3 send
// envelope-first packets, and both roots reach quorum on operator 3's packet.
func GloasProposerEnvelopeFirstOrder() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer envelope-first order",
		Documentation: testdoc.GloasProposerEnvelopeFirstOrderDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages: append(
			testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // pre-consensus + consensus
			[]*types.SignedSSVMessage{ // post-consensus: op1 (this operator) normal order; peers op2/op3 envelope-first
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerEnvelopeFirstMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerEnvelopeFirstMsgV(ks.Shares[3], 3, version))),
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
