package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RetainedPartialReconstructFailureReturned tests that replaying retained partials (SIP #94 §5) reports a
// quorum it completes but cannot reconstruct, as the live path does: op4's retained share has a bad
// signature, so the replayed quorum fails to reconstruct, the duty start returns that error, and the bad
// share is dropped. This operator's own share then re-crosses quorum and the preference submits.
func RetainedPartialReconstructFailureReturned() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// The peers' partials arrive before the duty starts, so the dispatcher retains them; op4's last.
	runner := testingutils.ProposerPreferencesRunner(ks)
	for _, msg := range []*types.PartialSignatureMessages{
		testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[2], 2),
		testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[3], 3),
		testingutils.PreConsensusProposerPreferencesWrongBeaconSigMsg(ks.Shares[4], 4),
	} {
		if err := runner.(*ssv.ProposerPreferencesRunner).ProcessPreConsensus(msg); err != nil {
			panic(err.Error())
		}
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences retained partial reconstruct failure returned",
		Documentation: testdoc.ProposerPreferencesRetainedPartialReconstructFailureReturnedDoc,
		Runner:        runner,
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		BeaconBroadcastedRoots: []string{
			// op2 + op3 + this operator's op1 share, once op4's bad share is dropped.
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas)),
		},
		ExpectedErrorCode: types.ReconstructSignatureErrorCode,
	}
}
