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

	state := State{
		DecidedValue: []byte{0xaa},
		Finished:     false,
		StartingDuty: &types.ValidatorDuty{Slot: 99},
	}
	err := json.Unmarshal([]byte(`{"Finished":true}`), &state)
	require.EqualError(t, err, "can't unmarshal BaseRunner.State.StartingDuty: expected one of ValidatorDuty, CommitteeDuty, or AggregatorCommitteeDuty")
	require.Equal(t, []byte{0xaa}, state.DecidedValue)
	require.False(t, state.Finished)
	require.Equal(t, phase0.Slot(99), state.StartingDuty.DutySlot())
}

func TestStateMarshalJSONRejectsMissingStartingDuty(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(&State{})
	require.EqualError(t, err, "json: error calling MarshalJSON for type *ssv.State: can't marshal BaseRunner.State.StartingDuty is nil")
}

func TestStateMarshalJSONRejectsUnsupportedStartingDuty(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(&State{StartingDuty: fakeDuty{}})
	require.EqualError(t, err, "json: error calling MarshalJSON for type *ssv.State: can't marshal BaseRunner.State.StartingDuty: expected ValidatorDuty, CommitteeDuty, or AggregatorCommitteeDuty, got: ssv.fakeDuty")
}

func TestStateUnmarshalJSONRejectsAmbiguousStartingDuty(t *testing.T) {
	t.Parallel()

	payload, err := json.Marshal(struct {
		Finished                bool                           `json:"Finished"`
		ValidatorDuty           *types.ValidatorDuty           `json:"ValidatorDuty,omitempty"`
		CommitteeDuty           *types.CommitteeDuty           `json:"CommitteeDuty,omitempty"`
		AggregatorCommitteeDuty *types.AggregatorCommitteeDuty `json:"AggregatorCommitteeDuty,omitempty"`
	}{
		Finished:                true,
		ValidatorDuty:           &types.ValidatorDuty{Slot: 1},
		CommitteeDuty:           &types.CommitteeDuty{Slot: 1},
		AggregatorCommitteeDuty: &types.AggregatorCommitteeDuty{Slot: 1},
	})
	require.NoError(t, err)

	state := State{
		DecidedValue: []byte{0xbb},
		Finished:     false,
		StartingDuty: &types.CommitteeDuty{Slot: 77},
	}
	err = json.Unmarshal(payload, &state)
	require.EqualError(t, err, "can't unmarshal BaseRunner.State.StartingDuty: payload contains more than one of ValidatorDuty, CommitteeDuty, and AggregatorCommitteeDuty")
	require.Equal(t, []byte{0xbb}, state.DecidedValue)
	require.False(t, state.Finished)
	require.Equal(t, phase0.Slot(77), state.StartingDuty.DutySlot())
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

	t.Run("aggregator committee duty", func(t *testing.T) {
		t.Parallel()

		payload, err := json.Marshal(struct {
			AggregatorCommitteeDuty *types.AggregatorCommitteeDuty `json:"AggregatorCommitteeDuty,omitempty"`
		}{
			AggregatorCommitteeDuty: &types.AggregatorCommitteeDuty{Slot: 56},
		})
		require.NoError(t, err)

		var state State
		err = json.Unmarshal(payload, &state)
		require.NoError(t, err)

		aggregatorCommitteeDuty, ok := state.StartingDuty.(*types.AggregatorCommitteeDuty)
		require.True(t, ok)
		require.EqualValues(t, 56, aggregatorCommitteeDuty.Slot)
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
		{
			name: "aggregator committee duty",
			duty: &types.AggregatorCommitteeDuty{Slot: 56},
			slot: 56,
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
			case *types.AggregatorCommitteeDuty:
				decodedDuty, ok := decoded.StartingDuty.(*types.AggregatorCommitteeDuty)
				require.True(t, ok)
				require.Equal(t, duty.Slot, decodedDuty.Slot)
			default:
				t.Fatalf("unsupported test duty type %T", tc.duty)
			}
		})
	}
}

func TestStateEncodeDecodeRoundTrip(t *testing.T) {
	t.Parallel()

	original := &State{
		DecidedValue: []byte{0x05, 0x06},
		Finished:     true,
		StartingDuty: &types.ValidatorDuty{Slot: 55},
	}

	payload, err := original.Encode()
	require.NoError(t, err)

	var decoded State
	err = decoded.Decode(payload)
	require.NoError(t, err)
	require.Equal(t, original.DecidedValue, decoded.DecidedValue)
	require.Equal(t, original.Finished, decoded.Finished)

	decodedDuty, ok := decoded.StartingDuty.(*types.ValidatorDuty)
	require.True(t, ok)
	require.Equal(t, phase0.Slot(55), decodedDuty.Slot)
}

// A dispatcher holding no retained partials omits PendingPartials from its JSON, so one rebuilt through Decode
// has a nil map and must still be able to retain the next partial it receives.
func TestProposerPreferencesRunnerRetainsPartialsAfterDecode(t *testing.T) {
	t.Parallel()

	const validatorIndex = phase0.ValidatorIndex(1)
	const slot = phase0.Slot(12)

	original := &ProposerPreferencesRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: types.RoleProposerPreferences,
			Share:          map[phase0.ValidatorIndex]*types.Share{validatorIndex: {ValidatorIndex: validatorIndex}},
		},
		BySlot:          map[phase0.Slot]*ProposerPreferencesSlotRunner{},
		PendingPartials: map[phase0.Slot][]*types.PartialSignatureMessages{},
	}
	payload, err := original.Encode()
	require.NoError(t, err)

	decoded := &ProposerPreferencesRunner{}
	require.NoError(t, decoded.Decode(payload))
	require.Nil(t, decoded.PendingPartials)

	partial := &types.PartialSignatureMessages{
		Type: types.ProposerPreferencesPartialSig,
		Slot: slot,
		Messages: []*types.PartialSignatureMessage{{
			PartialSignature: make([]byte, 96),
			Signer:           1,
			ValidatorIndex:   validatorIndex,
		}},
	}
	require.NoError(t, decoded.ProcessPreConsensus(partial))
	require.Len(t, decoded.PendingPartials[slot], 1)
}
