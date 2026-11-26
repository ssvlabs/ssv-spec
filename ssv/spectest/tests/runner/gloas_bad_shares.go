package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerBadShares tests an operator whose packet carries a bad block share and a bad §6 envelope share
// (SIP #94 §4/§6). Operator 3's packet brings both roots to quorum: the block fails to reconstruct and drops
// the share, and the envelope is held until the block submit is attempted. Operator 4's packet re-crosses the
// block, which submits. The envelope, at four shares, then fails to reconstruct and drops operator 3's share
// too. Three valid shares are left, a quorum that no later packet will cross again, so the envelope
// reconstructs from them at once and the reveal is published.
func GloasProposerBadShares() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer bad block and envelope shares",
		Documentation: testdoc.GloasProposerBadSharesDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages: append(
			testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // pre-consensus + consensus
			[]*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerBadSharesMsgV(ks.Shares[3], 3, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, version))),
			}...,
		),
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
		},
		BeaconBroadcastedRoots: testingutils.WithGloasEnvelopeBroadcast([]string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
		}, ks, version),
		OrderedBeaconBroadcastedRoots: true,
		ExpectedErrorCode:             types.ReconstructSignatureErrorCode,
	}
}
