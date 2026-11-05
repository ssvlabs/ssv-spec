package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RetainedPartialsCapIsolatesSigners tests that the retention cap counts per signer (SIP #94 §5/§7): a faulty
// peer floods committee-size × the per-signer budget of distinct auth partials before the duty starts, yet the
// honest peers' preference partials are still retained, so starting the duty reaches quorum and submits.
func RetainedPartialsCapIsolatesSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// Everything below arrives before the duty starts, so the dispatcher retains it: op2's flood, then ops 3
	// and 4's preference partials.
	runner := testingutils.ProposerPreferencesRunner(ks)
	msgs := append(
		distinctAuthPartials(ks, 2, len(ks.Committee())*(ssv.MaxProposerPreferencesDistinctRoots+ssv.MaxBuilderEntries)),
		testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[3], 3),
		testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[4], 4),
	)
	for _, msg := range msgs {
		if err := runner.(*ssv.ProposerPreferencesRunner).ProcessPreConsensus(msg); err != nil {
			panic(err.Error())
		}
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences retained partials cap isolates signers",
		Documentation: testdoc.ProposerPreferencesRetainedPartialsCapIsolatesSignersDoc,
		Runner:        runner,
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		BeaconBroadcastedRoots: []string{
			// retained op3 + op4 + this operator's op1 share = quorum, despite op2's flood.
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas)),
		},
	}
}
