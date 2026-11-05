package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DivergingDependentRoot tests honest convergence under a dependent-root split (SIP #94 §5): a peer
// whose beacon node computed the proposer duties against a different dependent root signs a different
// preference, so its partial fails the operator's expected-root check and the matching minority stays
// below quorum — divergence costs liveness, never a mixed signature.
func DivergingDependentRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences diverging dependent root",
		Documentation: testdoc.ProposerPreferencesDivergingDependentRootDoc,
		Runner:        testingutils.ProposerPreferencesRunner(ks),
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[2], 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesWrongRootMsg(ks.Shares[3], 3))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		ExpectedErrorCode: types.WrongSigningRootErrorCode,
	}
}
