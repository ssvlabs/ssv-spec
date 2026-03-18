package ssv

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
)

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
