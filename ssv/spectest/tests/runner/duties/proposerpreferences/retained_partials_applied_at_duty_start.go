package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RetainedPartialsAppliedAtDutyStart tests a late-starting operator (SIP #94 §5): three peers' preference
// partials arrive before its duty starts and are retained; starting the duty replays them, reaching quorum
// and submitting the preference.
func RetainedPartialsAppliedAtDutyStart() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// The peers' partials arrive before the duty starts, so the dispatcher retains them.
	runner := testingutils.ProposerPreferencesRunner(ks)
	for _, opID := range []types.OperatorID{2, 3, 4} {
		if err := runner.(*ssv.ProposerPreferencesRunner).ProcessPreConsensus(testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[opID], opID)); err != nil {
			panic(err.Error())
		}
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences retained partials applied at duty start",
		Documentation: testdoc.ProposerPreferencesRetainedPartialsAppliedAtDutyStartDoc,
		Runner:        runner,
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages:      []*types.SignedSSVMessage{}, // the duty start replays the retained partials
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas)),
		},
	}
}
