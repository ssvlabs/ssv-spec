package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerFailedEnvelopeSign tests that an operator whose signer fails on the §6 envelope root still
// broadcasts its block share: receivers do not require the envelope entry, so the failure costs this operator
// its envelope share only (SIP #94 §4). Its post-consensus packet goes out block-only and the signing error is
// returned. That block share then completes the block quorum (shares 1–3), while the envelope root reaches
// quorum on the other operators' shares alone (2, 3, 4), so this operator, which produced the decided block,
// still publishes the reveal.
func GloasProposerFailedEnvelopeSign() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer failed envelope sign",
		Documentation: testdoc.GloasProposerFailedEnvelopeSignDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages: append(
			testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // pre-consensus + consensus
			[]*types.SignedSSVMessage{ // post-consensus: this operator's own packet carries the block entry only
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerBlockOnlyMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, version))),
			}...,
		),
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			testingutils.PostConsensusProposerBlockOnlyMsgV(ks.Shares[1], 1, version),
		},
		BeaconBroadcastedRoots: testingutils.WithGloasEnvelopeBroadcast([]string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
		}, version),
		Signer:            &tests.SignerBehaviour{FailEnvelopeSigning: true},
		ExpectedErrorCode: types.ProposerEnvelopeSignFailedErrorCode,
	}
}
