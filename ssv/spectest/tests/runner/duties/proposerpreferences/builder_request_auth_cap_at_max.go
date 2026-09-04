package proposerpreferences

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthCapAtMax tests SSV's per-validator cap on configured builder entries (SIP #94 §5):
// with more than MaxBuilderEntries distinct entries, the round freezes and broadcasts partials for only
// the first MaxBuilderEntries. The assertion is the broadcast itself — no peer partials are fed, so no
// root reaches quorum — which pins the frozen set at exactly the cap.
func BuilderRequestAuthCapAtMax() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// More distinct entries than the cap; the round keeps the first MaxBuilderEntries in order.
	entries := make([]ssv.BuilderEntry, 0, ssv.MaxBuilderEntries+2)
	frozenData := make([][]byte, 0, ssv.MaxBuilderEntries)
	for i := 0; i < ssv.MaxBuilderEntries+2; i++ {
		data := []byte(fmt.Sprintf("cap-builder-auth-data-%02d", i))
		entries = append(entries, ssv.BuilderEntry{Data: data})
		if i < ssv.MaxBuilderEntries {
			frozenData = append(frozenData, data)
		}
	}

	runner := testingutils.ProposerPreferencesRunner(ks)
	runner.(*ssv.ProposerPreferencesRunner).BuilderEntries = entries

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences builder request auth cap at max",
		Documentation: testdoc.ProposerPreferencesBuilderRequestAuthCapAtMaxDoc,
		Runner:        runner,
		Duty:          testingutils.TestingProposerPreferencesDuty(),
		Messages:      []*types.SignedSSVMessage{},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1),            // preference partial, broadcast when starting a new duty
			testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, frozenData), // exactly MaxBuilderEntries auth partials
		},
		BeaconBroadcastedRoots: []string{},
	}
}
