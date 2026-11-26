package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RetainedAuthPartialsAppliedAtDutyStart is RetainedPartialsAppliedAtDutyStart for the request-auth round:
// three peers' two-entry auth packets arrive before the duty starts and are retained; starting the duty
// replays them, so both auth roots reach quorum and submit.
func RetainedAuthPartialsAppliedAtDutyStart() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()
	authData := [][]byte{data0, data1}

	// The peers' partials arrive before the duty starts, so the dispatcher retains them.
	var preDuty []*types.SignedSSVMessage
	for _, opID := range []types.OperatorID{2, 3, 4} {
		preDuty = append(preDuty, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[opID], opID, authData))))
	}

	return &tests.MsgProcessingSpecTest{
		Name:            "proposer preferences retained auth partials applied at duty start",
		Documentation:   testdoc.ProposerPreferencesRetainedAuthPartialsAppliedAtDutyStartDoc,
		Runner:          testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
		Duty:            testingutils.TestingProposerPreferencesDuty(),
		PreDutyMessages: preDuty,
		Messages:        []*types.SignedSSVMessage{}, // the duty start replays the retained partials
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1),          // preference partial, broadcast when starting a new duty
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData), // both auth partials, in one container
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data0, testingutils.TestingDutySlotGloas)),
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data1, testingutils.TestingDutySlotGloas)),
		},
	}
}
