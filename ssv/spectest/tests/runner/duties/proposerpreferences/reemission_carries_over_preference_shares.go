package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ReemissionCarriesOverPreferenceShares tests SIP #94 §5 preference-share carry-over across a same-slot
// re-emission — the pair to BuilderRequestAuthReemissionCarryOver on the auth side. The replacement
// sub-runner re-derives a byte-identical preference (same signing root, the dependent_root being unchanged),
// so the prior incarnation's collected shares carry over rather than resetting. The first incarnation holds
// two of the three shares needed for quorum; after the re-emission a single further share completes the
// quorum and submits — a reset that discarded the carried shares could not, since peers dedup the
// re-broadcast (SIP §7).
func ReemissionCarriesOverPreferenceShares() tests.SpecTest {
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
		Name:          "proposer preferences reemission carries over preference shares",
		Documentation: testdoc.ProposerPreferencesReemissionCarriesOverPreferenceSharesDoc,
		Runner:        runner,
		Duty:          duty, // the re-emission
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcast by the re-emission
		},
		BeaconBroadcastedRoots: []string{
			// carried op2 + op3 + the post-re-emission op1 share = quorum → the preference submits.
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, duty.Slot)),
		},
	}
}
