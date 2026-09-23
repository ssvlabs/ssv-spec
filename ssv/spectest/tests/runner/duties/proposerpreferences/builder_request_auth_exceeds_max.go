package proposerpreferences

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthExceedsMax tests the SIP #94 §7 receive-side bound: a RequestAuthPartialSig container
// carrying more than MaxBuilderEntries entries is rejected (RequestAuthWrongRootsCountErrorCode) rather than
// partially processed. The runner still freezes and broadcasts its own preference + auth partials on start;
// only the oversized peer packet is rejected.
func BuilderRequestAuthExceedsMax() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()
	ownAuthData := [][]byte{data0, data1}

	// One entry over the cap.
	oversized := make([][]byte, ssv.MaxBuilderEntries+1)
	for i := range oversized {
		oversized[i] = []byte(fmt.Sprintf("builder-auth-data-%d", i))
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences builder request auth exceeds max",
		Documentation: testdoc.ProposerPreferencesBuilderRequestAuthExceedsMaxDoc,
		Runner:        testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, oversized))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1),
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, ownAuthData),
		},
		ExpectedErrorCode: types.RequestAuthWrongRootsCountErrorCode,
	}
}
