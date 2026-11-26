package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RoundFailures tests that the duty's two rounds do not gate each other when one of them fails at the start
// (SIP #94 §5): the other is still signed and broadcast, both keep collecting, and the failure is returned.
//
//   - The request-auth container fails to broadcast: the preference partial still goes out, and the duty runs
//     with both rounds frozen, so the peers' partials complete each of them.
//   - The preference partial fails to broadcast: the request-auth container still goes out, and both rounds
//     complete in the same way.
//   - The preference round freezes nothing, because the dependent-root request fails: the request-auth round
//     still goes out and completes, while a peer's preference partial has nothing to be checked against.
func RoundFailures() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()
	authData := [][]byte{data0, data1}

	peerMsg := func(msg *types.PartialSignatureMessages) *types.SignedSSVMessage {
		return testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, msg))
	}
	// Three operators' auth containers, then their preference partials: a quorum for each root.
	var authMsgs, preferenceMsgs []*types.SignedSSVMessage
	for _, opID := range []types.OperatorID{1, 2, 3} {
		authMsgs = append(authMsgs, peerMsg(testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[opID], opID, authData)))
		preferenceMsgs = append(preferenceMsgs, peerMsg(testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[opID], opID)))
	}

	signedAuths := []string{
		testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data0, testingutils.TestingDutySlotGloas)),
		testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data1, testingutils.TestingDutySlotGloas)),
	}
	signedPreferences := testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, testingutils.TestingDutySlotGloas))

	return tests.NewMultiMsgProcessingSpecTest(
		"proposer preferences round failures",
		testdoc.ProposerPreferencesRoundFailuresDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:     "request auth broadcast fails",
				Runner:   testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
				Duty:     testingutils.TestingProposerPreferencesDuty(),
				Messages: append(append([]*types.SignedSSVMessage{}, authMsgs...), preferenceMsgs...),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // the auth container is not broadcast
				},
				BeaconBroadcastedRoots: append(append([]string{}, signedAuths...), signedPreferences),
				Network:                &tests.NetworkBehaviour{FailPartialSigBroadcast: []types.PartialSigMsgType{types.RequestAuthPartialSig}},
				ExpectedErrorCode:      types.ProposerPreferencesBroadcastFailedErrorCode,
			},
			{
				Name:     "preference broadcast fails",
				Runner:   testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
				Duty:     testingutils.TestingProposerPreferencesDuty(),
				Messages: append(append([]*types.SignedSSVMessage{}, authMsgs...), preferenceMsgs...),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData), // the preference partial is not broadcast
				},
				BeaconBroadcastedRoots: append(append([]string{}, signedAuths...), signedPreferences),
				Network:                &tests.NetworkBehaviour{FailPartialSigBroadcast: []types.PartialSigMsgType{types.ProposerPreferencesPartialSig}},
				ExpectedErrorCode:      types.ProposerPreferencesBroadcastFailedErrorCode,
			},
			{
				Name:   "preference round freezes nothing",
				Runner: testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
				Duty:   testingutils.TestingProposerPreferencesDuty(),
				// The last message is a peer's preference partial, rejected for the lack of a frozen preference.
				Messages: append(append([]*types.SignedSSVMessage{}, authMsgs...), preferenceMsgs[1]),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData),
				},
				BeaconBroadcastedRoots: signedAuths,
				BeaconNode:             &tests.BeaconNodeBehaviour{FailProposerDutiesDependentRoot: true},
				ExpectedErrorCode:      types.ProposerPreferencesNoPreferenceErrorCode,
			},
		},
		ks,
	)
}
