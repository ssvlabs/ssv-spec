package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ReemissionAfterSubmitStaysFinished tests a same-slot re-emission after the preference was already submitted
// (SIP #94 §5). The replacement sub-runner re-derives a byte-identical preference, so the prior incarnation's
// quorum carries over without firing again, and its completion carries with it: nothing is submitted a second
// time, the slot does not go back to running, and a late peer partial for the submitted preference is rejected
// as for any finished duty.
func ReemissionAfterSubmitStaysFinished() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// First incarnation: seeded as having collected a quorum (operators 1–3) and submitted the preference,
	// with those partials retained by the dispatcher as it retains every partial it receives.
	runner := testingutils.ProposerPreferencesRunner(ks)
	dispatcher := runner.(*ssv.ProposerPreferencesRunner)
	duty := testingutils.TestingProposerPreferencesDuty()
	firstSub := dispatcher.NewSlotRunner()
	firstSub.BaseRunner.State = ssv.NewRunnerState(ks.Threshold, duty)
	firstSub.ProposerPreferences = testingutils.TestingProposerPreferences(duty.Slot)
	for _, opID := range []types.OperatorID{1, 2, 3} {
		msg := testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[opID], opID)
		firstSub.BaseRunner.State.PreConsensusContainer.AddSignature(msg.Messages[0])
		dispatcher.PendingPartials[duty.Slot] = append(dispatcher.PendingPartials[duty.Slot], msg)
	}
	firstSub.BaseRunner.State.Finished = true
	dispatcher.BySlot[duty.Slot] = firstSub

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences reemission after submit stays finished",
		Documentation: testdoc.ProposerPreferencesReemissionAfterSubmitStaysFinishedDoc,
		Runner:        runner,
		Duty:          duty, // the re-emission
		Messages: []*types.SignedSSVMessage{
			// a late partial for the already-submitted preference
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[4], 4))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcast by the re-emission
		},
		// No beacon broadcast: the carried quorum does not submit again.
		ExpectedErrorCode: types.NoRunningDutyErrorCode,
	}
}
