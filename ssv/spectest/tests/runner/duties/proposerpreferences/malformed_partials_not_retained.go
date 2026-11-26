package proposerpreferences

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MalformedPartialsNotRetained tests that the dispatcher rejects a malformed packet before retaining any of
// its partials (SIP #94 §5/§7). Retained partials are replayed one at a time when the slot's duty starts, so a
// packet the live path rejects whole must not leave its well-formed partials behind to be collected then. Each
// packet arrives before its slot's duty has started; the post-state pins that nothing was retained.
func MalformedPartialsNotRetained() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	duty := testingutils.TestingProposerPreferencesDuty()

	// A matching preference partial with a diverging one behind it: more than the one partial a packet carries.
	twoPreferences := testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[2], 2)
	twoPreferences.Messages = append(twoPreferences.Messages,
		testingutils.PreConsensusProposerPreferencesWrongRootMsg(ks.Shares[2], 2).Messages...)

	// One request-auth partial over the cap.
	oversized := make([][]byte, ssv.MaxBuilderEntries+1)
	for i := range oversized {
		oversized[i] = []byte(fmt.Sprintf("builder-auth-data-%d", i))
	}

	// Two request-auth partials for the configured entries, the second for a validator this runner does not run.
	unknownValidator := testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, [][]byte{
		testingutils.TestingBuilderEntries[0].AuthData(),
		testingutils.TestingBuilderEntries[1].AuthData(),
	})
	unknownValidator.Messages[1].ValidatorIndex = testingutils.TestingWrongValidatorIndex

	return tests.NewMultiMsgProcessingSpecTest(
		"proposer preferences malformed partials not retained",
		testdoc.ProposerPreferencesMalformedPartialsNotRetainedDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "two preference partials",
				Runner: testingutils.ProposerPreferencesRunner(ks),
				Duty:   duty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, twoPreferences)),
				},
				OutputMessages:    []*types.PartialSignatureMessages{},
				DontStartDuty:     true,
				ExpectedErrorCode: types.WrongRootsCountErrorCode,
			},
			{
				Name:   "request auth partials exceed max",
				Runner: testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
				Duty:   duty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, oversized))),
				},
				OutputMessages:    []*types.PartialSignatureMessages{},
				DontStartDuty:     true,
				ExpectedErrorCode: types.RequestAuthWrongRootsCountErrorCode,
			},
			{
				Name:   "request auth partial for unknown validator",
				Runner: testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
				Duty:   duty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, unknownValidator)),
				},
				OutputMessages:    []*types.PartialSignatureMessages{},
				DontStartDuty:     true,
				ExpectedErrorCode: types.UnknownValidatorIndexErrorCode,
			},
		},
		ks,
	)
}
