package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthSharedDataDedup tests that entries sharing auth data share a single frozen
// builder-request-auth (SIP #94 §5): three entries, two carrying the same data, freeze two auths — not
// three — so the round broadcasts a two-partial container and reaches two per-root quorums.
func BuilderRequestAuthSharedDataDedup() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	dataA := []byte("shared-builder-auth-data")
	dataB := []byte("other-builder-auth-data")
	// Two of the three entries carry dataA, so the distinct-data set the round freezes is [A, B].
	entries := []ssv.BuilderEntry{
		{Data: dataA, URL: "https://builder-a.example"},
		{Data: dataB, URL: "https://builder-b.example"},
		{Data: dataA, URL: "https://builder-a-mirror.example"},
	}
	authData := [][]byte{dataA, dataB}

	runner := testingutils.ProposerPreferencesRunner(ks)
	runner.(*ssv.ProposerPreferencesRunner).BuilderEntries = entries

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences builder request auth shared data dedup",
		Documentation: testdoc.ProposerPreferencesBuilderRequestAuthSharedDataDedupDoc,
		Runner:        runner,
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
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData), // two auth partials, not three
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas)),
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, dataA, testingutils.TestingDutySlotGloas)),
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, dataB, testingutils.TestingDutySlotGloas)),
		},
	}
}
