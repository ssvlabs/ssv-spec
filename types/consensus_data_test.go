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

// TestConsensusDataVersionJSONLegacyNumeric pins the versionJSON codec's fallback for the legacy numeric Version
// form, which no generated vector carries (vectors round-trip the "gloas" string form).
func TestConsensusDataVersionJSONLegacyNumeric(t *testing.T) {
	var got ProposerConsensusData
	numeric := fmt.Sprintf(`{"Version":%d}`, uint64(gloas.DataVersionGloas))
	require.NoError(t, json.Unmarshal([]byte(numeric), &got))
	require.Equal(t, gloas.DataVersionGloas, got.Version)
}

// TestConsensusDataVersionJSON pins the versionJSON codec: a version go-eth2-client names is written as its fork
// string and read back in any case or from its number; any other version is written as its number, so that it
// reads back; a missing or null Version reads as 0.
func TestConsensusDataVersionJSON(t *testing.T) {
	const outOfEnum = spec.DataVersion(1 << 20)

	t.Run("marshal", func(t *testing.T) {
		for _, tc := range []struct {
			version spec.DataVersion
			want    string
		}{
			{gloas.DataVersionGloas, `"gloas"`},
			{spec.DataVersionElectra, `"electra"`},
			{spec.DataVersionUnknown, `0`},
			{outOfEnum, `1048576`},
		} {
			got, err := json.Marshal(versionJSON(tc.version))
			require.NoError(t, err)
			require.Equal(t, tc.want, string(got))
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		for _, tc := range []struct {
			input string
			want  spec.DataVersion
		}{
			{`"gloas"`, gloas.DataVersionGloas},
			{`"GLOAS"`, gloas.DataVersionGloas},
			{`"electra"`, spec.DataVersionElectra},
			{fmt.Sprint(uint64(gloas.DataVersionGloas)), gloas.DataVersionGloas},
			{`0`, spec.DataVersionUnknown},
			{`1048576`, outOfEnum},
			{`null`, spec.DataVersionUnknown},
		} {
			var got versionJSON
			require.NoError(t, json.Unmarshal([]byte(tc.input), &got), tc.input)
			require.Equal(t, tc.want, spec.DataVersion(got), tc.input)
		}

		var got versionJSON
		require.Error(t, json.Unmarshal([]byte(`"not-a-fork"`), &got))
	})

	// Both consensus-data types write their Version through the codec.
	t.Run("consensus data", func(t *testing.T) {
		for _, cd := range []any{
			&ProposerConsensusData{Version: gloas.DataVersionGloas},
			&AggregatorCommitteeConsensusData{Version: gloas.DataVersionGloas},
		} {
			byts, err := json.Marshal(cd)
			require.NoError(t, err)
			require.Contains(t, string(byts), `"Version":"gloas"`)
		}
	})
}
