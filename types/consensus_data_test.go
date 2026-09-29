package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

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
