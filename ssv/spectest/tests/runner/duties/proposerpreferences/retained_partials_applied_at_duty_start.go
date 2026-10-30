package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RetainedPartialsAppliedAtDutyStart tests a late-starting operator (SIP #94 §5): two peers' preference and
// request-auth partials arrive before its duty starts and are retained; starting the duty applies them, so
// its own share completes each quorum and the preference and both auths submit.
func RetainedPartialsAppliedAtDutyStart() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()
	authData := [][]byte{data0, data1}

	// The peers' partials arrive before the duty starts, so the dispatcher retains them.
	runner := testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks)
	for _, opID := range []types.OperatorID{2, 3} {
		for _, msg := range []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[opID], opID),
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[opID], opID, authData),
		} {
			if err := runner.(*ssv.ProposerPreferencesRunner).ProcessPreConsensus(msg); err != nil {
				panic(err.Error())
			}
		}
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences retained partials applied at duty start",
		Documentation: testdoc.ProposerPreferencesRetainedPartialsAppliedAtDutyStartDoc,
		Runner:        runner,
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1),          // preference partial, broadcast when starting a new duty
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData), // both auth partials, in one container
		},
		BeaconBroadcastedRoots: []string{
			// retained op2 + op3 + this operator's op1 share = quorum, for each of the three roots.
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas)),
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data0, testingutils.TestingDutySlotGloas)),
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data1, testingutils.TestingDutySlotGloas)),
		},
	}
}
