package ssv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TestProposerPreferencesRejectsUnexpectedPreConsensusType pins the dispatch guard in
// ProcessPreConsensus: the proposer-preferences runner multiplexes exactly two pre-consensus types
// (ProposerPreferencesPartialSig and RequestAuthPartialSig, SIP #94 §5), so any other type must be
// rejected outright rather than silently funneled into the preference path. This is a Go unit test
// rather than a spec vector because the guard is a fail-fast implementation detail — the
// cross-implementation contract stays the per-runner signing-domain/root check.
func TestProposerPreferencesRejectsUnexpectedPreConsensusType(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	runner := testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks)
	require.NoError(t, runner.StartNewDuty(testingutils.TestingProposerPreferencesDuty(), ks.Threshold))

	// A stray pre-consensus type addressed to this runner (here a PTC-attester partial) must be rejected.
	wrong := &types.PartialSignatureMessages{
		Type: types.PTCAttesterPartialSig,
		Slot: testingutils.TestingProposerPreferencesDuty().DutySlot(),
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: make([]byte, 96),
				Signer:           1,
				ValidatorIndex:   testingutils.TestingValidatorIndex,
			},
		},
	}

	err := runner.ProcessPreConsensus(wrong)
	require.Error(t, err)
	var specErr *types.Error
	require.ErrorAs(t, err, &specErr)
	require.Equal(t, types.ProposerPreferencesUnexpectedPartialSigTypeErrorCode, specErr.Code)
}
