package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthURLFallback tests the builder-request-auth data default (SIP #94 §5, builder-specs
// get_default_auth_data): an entry without Data signs the builder URL's lower-cased hostname.
func BuilderRequestAuthURLFallback() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	runner := testingutils.ProposerPreferencesRunner(ks)
	runner.(*ssv.ProposerPreferencesRunner).BuilderEntries = []ssv.BuilderEntry{{URL: "HTTPS://Builder.Example.com:443/bids?x=1"}}
	authData := [][]byte{[]byte("builder.example.com")}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences builder request auth url fallback",
		Documentation: testdoc.ProposerPreferencesBuilderRequestAuthURLFallbackDoc,
		Runner:        runner,
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, authData))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[3], 3, authData))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1),          // preference partial, broadcast when starting a new duty
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData), // the auth partial over the hostname
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, authData[0], testingutils.TestingDutySlotGloas)),
		},
	}
}
