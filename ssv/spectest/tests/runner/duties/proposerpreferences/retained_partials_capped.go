package proposerpreferences

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RetainedPartialsCapped tests the cap on a slot's retained partials (SIP #94 §5/§7): the committee size
// times the per-signer root budget (MaxProposerPreferencesDistinctRoots + MaxBuilderEntries). A peer sends
// one distinct auth partial more than that before the duty starts, and the excess one is not retained.
func RetainedPartialsCapped() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	limit := len(ks.Committee()) * (ssv.MaxProposerPreferencesDistinctRoots + ssv.MaxBuilderEntries)
	var messages []*types.SignedSSVMessage
	for sent := 0; sent <= limit; {
		var authData [][]byte
		for ; len(authData) < ssv.MaxBuilderEntries && sent <= limit; sent++ {
			authData = append(authData, []byte(fmt.Sprintf("builder-auth-token-%d", sent)))
		}
		messages = append(messages, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, authData))))
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences retained partials capped",
		Documentation: testdoc.ProposerPreferencesRetainedPartialsCappedDoc,
		Runner:        testingutils.ProposerPreferencesRunner(ks),
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages:      messages,
		DontStartDuty: true,
	}
}
