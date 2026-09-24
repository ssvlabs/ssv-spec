package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthReemissionCarryOver tests SIP #94 §5 auth-share carry-over across a same-slot
// re-emission — the preference-round counterpart is ReemissionCarriesOverPreferenceShares. Auth roots carry
// no dependent_root, so a re-emission always re-freezes byte-identical roots and already-collected auth shares
// carry over rather than restarting from zero (peers dedup the re-broadcast, SIP §7). The first incarnation holds two of the
// three data0 shares needed for quorum; after the re-emission a single further share completes the quorum
// and submits — a reset that discarded the carried shares could not.
func BuilderRequestAuthReemissionCarryOver() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	runner := testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks)
	duty := testingutils.TestingProposerPreferencesDuty()
	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()
	authData := [][]byte{data0, data1}

	// First incarnation: seeded as already emitted, with the auth roots frozen and two data0 shares
	// aggregated (one below quorum).
	firstSub := runner.(*ssv.ProposerPreferencesRunner).NewSlotRunner()
	firstSub.BaseRunner.State = ssv.NewRunnerState(ks.Threshold, duty)
	firstSub.BuilderRequestAuths = []*gloas.BuilderRequestAuth{
		{Data: data0, Slot: duty.Slot},
		{Data: data1, Slot: duty.Slot},
	}
	for _, opID := range []types.OperatorID{2, 3} {
		if err := firstSub.ProcessPreConsensus(testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[opID], opID, [][]byte{data0})); err != nil {
			panic(err.Error())
		}
	}
	runner.(*ssv.ProposerPreferencesRunner).BySlot[duty.Slot] = firstSub

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences builder request auth reemission carry over",
		Documentation: testdoc.ProposerPreferencesBuilderRequestAuthReemissionCarryOverDoc,
		Runner:        runner,
		Duty:          duty, // the re-emission
		Messages: []*types.SignedSSVMessage{
			// The re-emission carries op2 and op3's data0 shares over; op1's data0 share completes quorum.
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, [][]byte{data0}))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1),          // preference partial, broadcast by the re-emission
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData), // both auth partials, broadcast by the re-emission
		},
		BeaconBroadcastedRoots: []string{
			// carried op2 + op3 + the post-re-emission op1 share = quorum for data0.
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data0, testingutils.TestingDutySlotGloas)),
		},
	}
}
