package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoPreference tests a duty whose preference round froze nothing (SIP #94 §5): the beacon node fails the
// proposer-duties dependent-root request, so the operator derives no preference and signs and broadcasts nothing.
// With no frozen preference an incoming peer partial is rejected (ProposerPreferencesNoPreferenceErrorCode). It
// stays retained, for a later re-emission that does freeze one. The duty stays running.
func NoPreference() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences no preference",
		Documentation: testdoc.ProposerPreferencesNoPreferenceDoc,
		Runner:        testingutils.ProposerPreferencesRunner(ks),
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[2], 2))),
		},
		OutputMessages:    []*types.PartialSignatureMessages{},
		BeaconNode:        &tests.BeaconNodeBehaviour{FailProposerDutiesDependentRoot: true},
		ExpectedErrorCode: types.ProposerPreferencesNoPreferenceErrorCode,
	}
}
