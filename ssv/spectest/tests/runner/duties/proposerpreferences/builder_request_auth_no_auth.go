package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthNoAuth tests SIP #94 §5: an operator that froze no auth roots — no configured builder
// entries — rejects an incoming RequestAuthPartialSig with RequestAuthNoAuthErrorCode rather than collecting
// against an empty set. The preference round still runs (it needs no builder entries), so the runner
// broadcasts only its preference partial on start.
func BuilderRequestAuthNoAuth() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences builder request auth no auth",
		Documentation: testdoc.ProposerPreferencesBuilderRequestAuthNoAuthDoc,
		Runner:        testingutils.ProposerPreferencesRunner(ks), // no builder entries configured
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, [][]byte{data0, data1}))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // preference only; no auth round without entries
		},
		ExpectedErrorCode: types.RequestAuthNoAuthErrorCode,
	}
}
