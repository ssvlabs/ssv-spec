package ssv_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TestProposerPreferencesAuthBadShareDoesNotStrandOtherRoots pins SIP #94 §5 per-root isolation on the
// builder-request-auth reconstruction/submit side: when one packet pushes several auth roots over quorum at
// once and one root's share is invalid, that root fails to reconstruct but the others must still submit — a
// bad share for one builder must never strand another's auth. This is a Go unit test rather than a spec
// vector because the reconstruction failure surfaces as an untyped BLS error, which the MsgProcessingSpecTest
// ExpectedErrorCode framework (typed errors only) cannot match.
func TestProposerPreferencesAuthBadShareDoesNotStrandOtherRoots(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	runner := testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks)
	bn := runner.GetBeaconNode().(*testingutils.TestingBeaconNode)

	require.NoError(t, runner.StartNewDuty(testingutils.TestingProposerPreferencesDuty(), ks.Threshold))

	data0 := testingutils.TestingBuilderEntries[0].AuthData()
	data1 := testingutils.TestingBuilderEntries[1].AuthData()
	authData := [][]byte{data0, data1}

	// Two valid shares (op1 self, op2) for each root — below the quorum of three, so nothing submits yet.
	require.NoError(t, runner.ProcessPreConsensus(testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, authData)))
	require.NoError(t, runner.ProcessPreConsensus(testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[2], 2, authData)))
	require.Empty(t, bn.BroadcastedRoots, "no auth submitted before quorum")

	// op3's packet pushes both roots to quorum at once, but its data0 share is invalid. data0 fails
	// reconstruction; data1 must still submit (the fix), rather than an early return stranding it.
	err := runner.ProcessPreConsensus(testingutils.PreConsensusBuilderRequestAuthBadFirstShareMsg(ks.Shares[3], 3, authData))
	require.Error(t, err, "data0's invalid signature is surfaced")
	require.Contains(t, err.Error(), "invalid signatures")

	root0 := testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data0, testingutils.TestingDutySlotGloas))
	root1 := testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, data1, testingutils.TestingDutySlotGloas))
	var submitted []string
	for _, r := range bn.BroadcastedRoots {
		submitted = append(submitted, hex.EncodeToString(r[:]))
	}
	require.Contains(t, submitted, root1, "data1 must be submitted despite data0's bad share")
	require.NotContains(t, submitted, root0, "data0 must not be submitted (its share was invalid)")
}
