package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RetainedDivergingPartialDropped tests that retained partials (SIP #94 §5) are checked one at a time when
// the duty starts: the first, from a peer whose preference diverged, fails the expected-root check and is
// dropped with its error returned, while the two matching ones behind it still apply and reach quorum.
func RetainedDivergingPartialDropped() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// The peers' partials arrive before the duty starts, so the dispatcher retains them; the diverging one first.
	runner := testingutils.ProposerPreferencesRunner(ks)
	for _, msg := range []*types.PartialSignatureMessages{
		testingutils.PreConsensusProposerPreferencesWrongRootMsg(ks.Shares[2], 2),
		testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[3], 3),
		testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[4], 4),
	} {
		if err := runner.(*ssv.ProposerPreferencesRunner).ProcessPreConsensus(msg); err != nil {
			panic(err.Error())
		}
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences retained diverging partial dropped",
		Documentation: testdoc.ProposerPreferencesRetainedDivergingPartialDroppedDoc,
		Runner:        runner,
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		BeaconBroadcastedRoots: []string{
			// retained op3 + op4 + this operator's op1 share = quorum; op2's diverging share was dropped.
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas)),
		},
		ExpectedErrorCode: types.WrongSigningRootErrorCode,
	}
}
