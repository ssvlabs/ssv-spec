package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthHappyFlow tests the §5 builder-request-auth round riding the proposer-preferences
// duty (SIP #94 §5 builder-request-auth extension). Two distinct-data builder entries are configured, so
// executing the duty broadcasts the preference partial plus one RequestAuthPartialSig container carrying
// both auth partials, and each of the three roots reaches quorum and submits.
func BuilderRequestAuthHappyFlow() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()
	authData := [][]byte{data0, data1}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences builder request auth happy flow",
		Documentation: testdoc.ProposerPreferencesBuilderRequestAuthHappyFlowDoc,
		Runner:        testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[2], 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[3], 3))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, authData))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[3], 3, authData))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1),          // preference partial, broadcast when starting a new duty
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData), // both auth partials, in one container
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas)),
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data0, testingutils.TestingDutySlotGloas)),
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data1, testingutils.TestingDutySlotGloas)),
		},
	}
}
