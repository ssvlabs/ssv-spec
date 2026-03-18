package ssv

import (
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
)

type fakeDuty struct{}

func (fakeDuty) DutySlot() phase0.Slot {
	return 0
}

func (fakeDuty) RunnerRole() types.RunnerRole {
	return types.RoleUnknown
}

func TestStateUnmarshalJSONRejectsMissingStartingDuty(t *testing.T) {
	t.Parallel()

	var state State
	err := json.Unmarshal([]byte(`{"Finished":true}`), &state)
	require.EqualError(t, err, "can't unmarshal BaseRunner.State.StartingDuty: expected ValidatorDuty or CommitteeDuty")
}

func TestStateMarshalJSONRejectsMissingStartingDuty(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(&State{})
	require.EqualError(t, err, "json: error calling MarshalJSON for type *ssv.State: can't marshal BaseRunner.State.StartingDuty: expected ValidatorDuty or CommitteeDuty")
}

func TestStateMarshalJSONRejectsUnsupportedStartingDuty(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(&State{StartingDuty: fakeDuty{}})
	require.EqualError(t, err, "json: error calling MarshalJSON for type *ssv.State: can't marshal BaseRunner.State.StartingDuty: expected ValidatorDuty or CommitteeDuty")
}

func TestStateUnmarshalJSONRejectsAmbiguousStartingDuty(t *testing.T) {
	t.Parallel()

	payload, err := json.Marshal(struct {
		ValidatorDuty *types.ValidatorDuty `json:"ValidatorDuty,omitempty"`
		CommitteeDuty *types.CommitteeDuty `json:"CommitteeDuty,omitempty"`
	}{
		ValidatorDuty: &types.ValidatorDuty{Slot: 1},
		CommitteeDuty: &types.CommitteeDuty{Slot: 1},
	})
	require.NoError(t, err)

	var state State
	err = json.Unmarshal(payload, &state)
	require.EqualError(t, err, "can't unmarshal BaseRunner.State.StartingDuty: payload contains both ValidatorDuty and CommitteeDuty")
}

func TestStateUnmarshalJSONAcceptsKnownDutyTypes(t *testing.T) {
	t.Parallel()

	t.Run("validator duty", func(t *testing.T) {
		t.Parallel()

		payload, err := json.Marshal(struct {
			ValidatorDuty *types.ValidatorDuty `json:"ValidatorDuty,omitempty"`
		}{
			ValidatorDuty: &types.ValidatorDuty{Slot: 12},
		})
		require.NoError(t, err)

		var state State
		err = json.Unmarshal(payload, &state)
		require.NoError(t, err)

		validatorDuty, ok := state.StartingDuty.(*types.ValidatorDuty)
		require.True(t, ok)
		require.EqualValues(t, 12, validatorDuty.Slot)
	})

	t.Run("committee duty", func(t *testing.T) {
		t.Parallel()

		payload, err := json.Marshal(struct {
			CommitteeDuty *types.CommitteeDuty `json:"CommitteeDuty,omitempty"`
		}{
			CommitteeDuty: &types.CommitteeDuty{Slot: 34},
		})
		require.NoError(t, err)

		var state State
		err = json.Unmarshal(payload, &state)
		require.NoError(t, err)

		committeeDuty, ok := state.StartingDuty.(*types.CommitteeDuty)
		require.True(t, ok)
		require.EqualValues(t, 34, committeeDuty.Slot)
	})
}

func TestStateJSONRoundTripPreservesKnownDutyTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		duty types.Duty
		slot phase0.Slot
	}{
		{
			name: "validator duty",
			duty: &types.ValidatorDuty{Slot: 12},
			slot: 12,
		},
		{
			name: "committee duty",
			duty: &types.CommitteeDuty{Slot: 34},
			slot: 34,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			original := &State{
				DecidedValue: []byte{0x01, 0x02, 0x03},
				Finished:     true,
				StartingDuty: tc.duty,
			}

			payload, err := json.Marshal(original)
			require.NoError(t, err)

			var decoded State
			err = json.Unmarshal(payload, &decoded)
			require.NoError(t, err)
			require.Equal(t, original.DecidedValue, decoded.DecidedValue)
			require.Equal(t, original.Finished, decoded.Finished)
			require.Equal(t, tc.slot, decoded.StartingDuty.DutySlot())

			switch duty := tc.duty.(type) {
			case *types.ValidatorDuty:
				decodedDuty, ok := decoded.StartingDuty.(*types.ValidatorDuty)
				require.True(t, ok)
				require.Equal(t, duty.Slot, decodedDuty.Slot)
			case *types.CommitteeDuty:
				decodedDuty, ok := decoded.StartingDuty.(*types.CommitteeDuty)
				require.True(t, ok)
				require.Equal(t, duty.Slot, decodedDuty.Slot)
			default:
				t.Fatalf("unsupported test duty type %T", tc.duty)
			}
		})
	}
}
