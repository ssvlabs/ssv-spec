package proposerpreferences

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RetainedPartialsCapped tests the per-signer cap on a slot's retained partials (SIP #94 §5/§7): the §7
// per-signer root budget, MaxProposerPreferencesDistinctRoots + MaxBuilderEntries. A peer sends one distinct
// auth partial more than that before the duty starts, and the excess one is not retained.
func RetainedPartialsCapped() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	var messages []*types.SignedSSVMessage
	for _, msg := range distinctAuthPartials(ks, 2, ssv.MaxProposerPreferencesDistinctRoots+ssv.MaxBuilderEntries+1) {
		messages = append(messages, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, msg)))
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

// distinctAuthPartials returns count distinct request-auth partials from one signer, packed MaxBuilderEntries
// to a message.
func distinctAuthPartials(ks *testingutils.TestKeySet, signer types.OperatorID, count int) []*types.PartialSignatureMessages {
	var msgs []*types.PartialSignatureMessages
	for sent := 0; sent < count; {
		var authData [][]byte
		for ; len(authData) < ssv.MaxBuilderEntries && sent < count; sent++ {
			authData = append(authData, []byte(fmt.Sprintf("builder-auth-token-%d", sent)))
		}
		msgs = append(msgs, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[signer], signer, authData))
	}
	return msgs
}
