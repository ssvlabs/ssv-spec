package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ConcurrentLookaheadSlots tests two concurrently-active proposal slots (SIP #94 §5): a validator can
// hold several lookahead proposal slots at once, each an independent per-slot flow, so partials for
// both interleave and both reach quorum and submit — the property a single-state runner cannot
// provide (anchor's "slot-advance exemption").
func ConcurrentLookaheadSlots() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// The first proposal slot's flow is seeded as already emitted (frozen preference, running state);
	// the harness starts the second slot's duty, so both are active.
	runner := testingutils.ProposerPreferencesRunner(ks)
	firstDuty := testingutils.TestingProposerPreferencesDuty()
	firstSub := runner.(*ssv.ProposerPreferencesRunner).NewSlotRunner()
	firstSub.BaseRunner.State = ssv.NewRunnerState(ks.Threshold, firstDuty)
	firstSub.ProposerPreferences = testingutils.TestingProposerPreferences(firstDuty.Slot)
	runner.(*ssv.ProposerPreferencesRunner).BySlot[firstDuty.Slot] = firstSub

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences concurrent lookahead slots",
		Documentation: testdoc.ProposerPreferencesConcurrentLookaheadSlotsDoc,
		Runner:        runner,
		Duty:          testingutils.TestingProposerPreferencesSecondDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesSecondSlotMsg(ks.Shares[1], 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[2], 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesSecondSlotMsg(ks.Shares[2], 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[3], 3))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesSecondSlotMsg(ks.Shares[3], 3))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesSecondSlotMsg(ks.Shares[1], 1), // broadcast when the second slot's duty starts
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas)),
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas+1)),
		},
	}
}
