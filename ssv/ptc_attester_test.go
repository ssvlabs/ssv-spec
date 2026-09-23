package ssv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TestPTCAttesterRejectedDutyKeepsObservation covers the SIP #94 §3 observation lifecycle: a running duty
// freezes its payload-attestation observation, and a later duplicate/past-slot duty that
// ShouldProcessNonBeaconDuty rejects must not erase it. The clear lives at the top of executeDuty (reached
// only after a duty is accepted), so a rejected StartNewDuty leaves the observation intact. Before the fix
// the clear ran in StartNewDuty, so a rejected duplicate nil-ed a still-running duty's observation.
func TestPTCAttesterRejectedDutyKeepsObservation(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	r := testingutils.PTCAttesterRunner(ks).(*ssv.PTCAttesterRunner)

	// Running duty: executeDuty observes a block for the slot and freezes it.
	duty := testingutils.TestingPTCAttesterDuty()
	require.NoError(t, r.StartNewDuty(duty, ks.Threshold))
	frozen := r.PayloadAttestationData
	require.NotNil(t, frozen)
	require.Equal(t, testingutils.TestingPayloadAttestationData(duty.Slot), frozen)

	// A duplicate duty for the same slot is rejected by ShouldProcessNonBeaconDuty...
	err := r.StartNewDuty(testingutils.TestingPTCAttesterDuty(), ks.Threshold)
	require.Error(t, err)
	var typedErr *types.Error
	require.ErrorAs(t, err, &typedErr)
	require.Equal(t, types.DutyAlreadyPassedErrorCode, typedErr.Code)

	// ...and leaves the running duty's frozen observation untouched (regression: it was nil-ed).
	require.Same(t, frozen, r.PayloadAttestationData)
	require.Equal(t, testingutils.TestingPayloadAttestationData(duty.Slot), r.PayloadAttestationData)
}

// TestPTCAttesterRejectsWrongSlotData covers the SIP #94 §3 slot pin: if the beacon node answers with
// payload attestation data whose Slot differs from the duty slot, executeDuty rejects it
// (PTCAttesterWrongSlotErrorCode) rather than signing/broadcasting data for another slot under the duty
// slot's domain, and freezes nothing.
func TestPTCAttesterRejectsWrongSlotData(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	r := testingutils.PTCAttesterRunner(ks).(*ssv.PTCAttesterRunner)
	r.GetBeaconNode().(*testingutils.TestingBeaconNode).SetWrongPayloadAttestationSlot(true)

	err := r.StartNewDuty(testingutils.TestingPTCAttesterDuty(), ks.Threshold)
	require.Error(t, err)
	var typedErr *types.Error
	require.ErrorAs(t, err, &typedErr)
	require.Equal(t, types.PTCAttesterWrongSlotErrorCode, typedErr.Code)
	require.Nil(t, r.PayloadAttestationData) // nothing frozen on rejection
}
