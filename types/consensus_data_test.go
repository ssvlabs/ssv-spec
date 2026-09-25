package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// TestProposerConsensusDataJSONFieldsInSync guards the explicit ProposerConsensusData MarshalJSON/UnmarshalJSON
// overrides, which list the fields by hand to keep the JSON key order (Duty, Version, DataSSZ). A field added
// to the struct would otherwise silently vanish from JSON in both directions; if this fails, add the new
// field to both overrides in consensus_data.go too.
func TestProposerConsensusDataJSONFieldsInSync(t *testing.T) {
	require.Equal(t, 3, reflect.TypeOf(ProposerConsensusData{}).NumField(),
		"ProposerConsensusData gained a field — update its MarshalJSON/UnmarshalJSON overrides")
}

// TestConsensusDataVersionJSONRoundTrip pins the Gloas-safe Version JSON codec (versionJSON): a Gloas-stamped
// ProposerConsensusData / AggregatorCommitteeConsensusData marshals with Version "gloas" — spec.DataVersion's
// own MarshalJSON panics on the out-of-enum placeholder — and decodes back to DataVersionGloas (also from the
// legacy numeric form), while a pre-Gloas version keeps its upstream fork string. No generated fixture carries
// a Gloas-stamped Version, so this is the only thing exercising that branch.
func TestConsensusDataVersionJSONRoundTrip(t *testing.T) {
	t.Run("proposer gloas round-trips as \"gloas\"", func(t *testing.T) {
		cd := &ProposerConsensusData{Version: gloas.DataVersionGloas, DataSSZ: []byte{1, 2, 3}}
		b, err := json.Marshal(cd)
		require.NoError(t, err)
		require.Contains(t, string(b), `"Version":"gloas"`)

		var got ProposerConsensusData
		require.NoError(t, json.Unmarshal(b, &got))
		require.Equal(t, gloas.DataVersionGloas, got.Version)
		require.Equal(t, cd.DataSSZ, got.DataSSZ)
	})

	t.Run("proposer decodes the legacy numeric version", func(t *testing.T) {
		var got ProposerConsensusData
		numeric := fmt.Sprintf(`{"Version":%d}`, uint64(gloas.DataVersionGloas))
		require.NoError(t, json.Unmarshal([]byte(numeric), &got))
		require.Equal(t, gloas.DataVersionGloas, got.Version)
	})

	t.Run("proposer pre-gloas keeps the fork string", func(t *testing.T) {
		cd := &ProposerConsensusData{Version: spec.DataVersionElectra}
		b, err := json.Marshal(cd)
		require.NoError(t, err)
		require.Contains(t, string(b), `"Version":"electra"`)

		var got ProposerConsensusData
		require.NoError(t, json.Unmarshal(b, &got))
		require.Equal(t, spec.DataVersionElectra, got.Version)
	})

	t.Run("aggregator gloas round-trips as \"gloas\"", func(t *testing.T) {
		a := &AggregatorCommitteeConsensusData{Version: gloas.DataVersionGloas}
		b, err := json.Marshal(a)
		require.NoError(t, err)
		require.Contains(t, string(b), `"Version":"gloas"`)

		var got AggregatorCommitteeConsensusData
		require.NoError(t, json.Unmarshal(b, &got))
		require.Equal(t, gloas.DataVersionGloas, got.Version)
	})
}
