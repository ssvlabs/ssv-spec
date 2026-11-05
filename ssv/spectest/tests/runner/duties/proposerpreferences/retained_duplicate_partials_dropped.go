package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RetainedDuplicatePartialsDropped tests that the dispatcher retains each (signer, signing root) once
// (SIP #94 §5): a peer's preference packet and two-entry auth packet, each delivered twice before the duty
// starts, leave one preference and two auth partials retained.
func RetainedDuplicatePartialsDropped() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	authData := [][]byte{testingutils.TestingBuilderEntries[0].AuthData(), testingutils.TestingBuilderEntries[1].AuthData()}
	preference := testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[2], 2)))
	auth := testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, authData)))

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences retained duplicate partials dropped",
		Documentation: testdoc.ProposerPreferencesRetainedDuplicatePartialsDroppedDoc,
		Runner:        testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages:      []*types.SignedSSVMessage{preference, auth, preference, auth},
		DontStartDuty: true,
	}
}
