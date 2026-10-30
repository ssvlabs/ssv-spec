package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ReemissionAppliesPartialsPriorPreferenceRejected tests a dependent-root change (SIP #94 §5): the slot's
// first incarnation froze a preference on a stale dependent root, so it rejects the peers' partials over
// the current one, but the dispatcher still retains them. The re-emission freezes the current preference
// and replays them, reaching quorum and submitting.
func ReemissionAppliesPartialsPriorPreferenceRejected() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// First incarnation: seeded as already emitted, on a stale dependent root.
	runner := testingutils.ProposerPreferencesRunner(ks)
	duty := testingutils.TestingProposerPreferencesDuty()
	staleSub := runner.(*ssv.ProposerPreferencesRunner).NewSlotRunner()
	staleSub.BaseRunner.State = ssv.NewRunnerState(ks.Threshold, duty)
	staleSub.ProposerPreferences = testingutils.TestingProposerPreferences(duty.Slot)
	staleSub.ProposerPreferences.DependentRoot = testingutils.TestingBlockRoot
	runner.(*ssv.ProposerPreferencesRunner).BySlot[duty.Slot] = staleSub

	// The peers' partials over the current dependent root: the first incarnation rejects them, the
	// dispatcher retains them.
	for _, opID := range []types.OperatorID{2, 3, 4} {
		if err := runner.(*ssv.ProposerPreferencesRunner).ProcessPreConsensus(testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[opID], opID)); err == nil {
			panic("the stale incarnation accepted a partial over the current dependent root")
		}
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences reemission applies partials prior preference rejected",
		Documentation: testdoc.ProposerPreferencesReemissionAppliesPartialsPriorPreferenceRejectedDoc,
		Runner:        runner,
		Duty:          duty,                        // the re-emission, on the current dependent root
		Messages:      []*types.SignedSSVMessage{}, // the duty start replays the retained partials
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcast by the re-emission
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, duty.Slot)),
		},
	}
}
