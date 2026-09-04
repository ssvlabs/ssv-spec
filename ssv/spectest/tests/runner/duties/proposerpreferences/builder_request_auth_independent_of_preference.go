package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthIndependentOfPreference tests that the §5 builder-request-auth round keeps collecting
// after the preference round has finished (SIP #94 §5 — neither round gates the other). The slot is seeded
// with a finished preference round and two frozen auths; auth partials arriving afterward still reach
// per-root quorum and submit.
func BuilderRequestAuthIndependentOfPreference() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()
	authData := [][]byte{data0, data1}

	// Seed a slot whose preference round has already finished; only the auth round is still open.
	runner := testingutils.ProposerPreferencesRunner(ks)
	duty := testingutils.TestingProposerPreferencesDuty()
	sub := runner.(*ssv.ProposerPreferencesRunner).NewSlotRunner()
	sub.BaseRunner.State = ssv.NewRunnerState(ks.Threshold, duty)
	sub.BaseRunner.State.Finished = true
	sub.ProposerPreferences = testingutils.TestingProposerPreferences(duty.Slot)
	sub.BuilderRequestAuths = []*gloas.BuilderRequestAuth{
		testingutils.TestingBuilderRequestAuth(data0, duty.Slot),
		testingutils.TestingBuilderRequestAuth(data1, duty.Slot),
	}
	runner.(*ssv.ProposerPreferencesRunner).BySlot[duty.Slot] = sub

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences builder request auth independent of preference",
		Documentation: testdoc.ProposerPreferencesBuilderRequestAuthIndependentOfPreferenceDoc,
		Runner:        runner,
		Duty:          duty,
		DontStartDuty: true, // the preference round is already finished; only auth partials arrive
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, authData))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[3], 3, authData))),
		},
		OutputMessages: []*types.PartialSignatureMessages{},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data0, testingutils.TestingDutySlotGloas)),
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data1, testingutils.TestingDutySlotGloas)),
		},
	}
}
