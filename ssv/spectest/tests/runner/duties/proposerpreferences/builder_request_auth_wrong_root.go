package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthWrongRoot tests that a peer whose configured builder entries disagree is rejected
// (SIP #94 §5): its RequestAuthPartialSig container has the right count but one partial signs a divergent
// auth data, so it fails the operator's expected-root check rather than mixing into a quorum — divergence
// costs liveness, never a mixed signature.
func BuilderRequestAuthWrongRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()
	authData := [][]byte{data0, data1}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences builder request auth wrong root",
		Documentation: testdoc.ProposerPreferencesBuilderRequestAuthWrongRootDoc,
		Runner:        testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthWrongRootMsg(ks.Shares[2], 2, data0))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1),          // preference partial, broadcast when starting a new duty
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData), // both auth partials, in one container
		},
		BeaconBroadcastedRoots: []string{},
		ExpectedErrorCode:      types.WrongSigningRootErrorCode,
	}
}
