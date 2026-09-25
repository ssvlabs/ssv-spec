package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthWrongRoot tests SIP #94 §5/§7 per-builder isolation on the receive side: a peer whose
// multi-entry RequestAuthPartialSig carries one divergent (non-frozen) auth entry
// has that entry ignored, not the whole packet rejected, so the entries that do match a frozen root still
// collect and reach quorum. Here op2 diverges on the second entry — data0 (carried by all three) still
// reaches quorum and submits, while data1 (which op2 replaced with the divergent entry) falls one share
// short. The divergence costs that builder, never the slot.
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
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData))),
			// op2's second entry signs a divergent (non-frozen) auth data; it is ignored, data0 kept.
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthWrongRootMsg(ks.Shares[2], 2, data0))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[3], 3, authData))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1),          // preference partial, broadcast when starting a new duty
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData), // both auth partials, in one container
		},
		BeaconBroadcastedRoots: []string{
			// data0 reaches quorum (op1, op2, op3); data1 falls short (op1, op3) because op2 diverged.
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data0, testingutils.TestingDutySlotGloas)),
		},
	}
}
