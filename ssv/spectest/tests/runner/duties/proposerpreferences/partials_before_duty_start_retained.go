package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialsBeforeDutyStartRetained tests that a peer's preference and request-auth partials for a proposal
// slot whose duty has not started here yet are retained for that slot, not rejected (SIP #94 §5).
func PartialsBeforeDutyStartRetained() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	authData := [][]byte{testingutils.TestingBuilderEntries[0].AuthData(), testingutils.TestingBuilderEntries[1].AuthData()}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences partials before duty start retained",
		Documentation: testdoc.ProposerPreferencesPartialsBeforeDutyStartRetainedDoc,
		Runner:        testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[2], 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, authData))),
		},
		DontStartDuty: true,
	}
}
