package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ReemissionReplacesSlot tests that re-emitting a proposal slot's duty replaces the slot's flow
// (SIP #94 §5): the replacement freezes a freshly derived preference and starts a fresh signature
// container. The prior incarnation had two of the three partials needed for quorum; had the
// replacement carried its container over, the one post-re-emission partial would complete a quorum
// and submit — so the asserted absence of a submission is what pins the replacement.
func ReemissionReplacesSlot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// First incarnation: seeded as already emitted, with two partials aggregated (one below quorum).
	runner := testingutils.ProposerPreferencesRunner(ks)
	duty := testingutils.TestingProposerPreferencesDuty()
	firstSub := runner.(*ssv.ProposerPreferencesRunner).NewSlotRunner()
	firstSub.BaseRunner.State = ssv.NewRunnerState(ks.Threshold, duty)
	firstSub.ProposerPreferences = testingutils.TestingProposerPreferences(duty.Slot)
	for _, opID := range []types.OperatorID{2, 3} {
		if err := firstSub.ProcessPreConsensus(testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[opID], opID)); err != nil {
			panic(err.Error())
		}
	}
	runner.(*ssv.ProposerPreferencesRunner).BySlot[duty.Slot] = firstSub

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences reemission replaces slot",
		Documentation: testdoc.ProposerPreferencesReemissionReplacesSlotDoc,
		Runner:        runner,
		Duty:          duty, // the re-emission
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcast by the re-emission
		},
		BeaconBroadcastedRoots: []string{},
	}
}
