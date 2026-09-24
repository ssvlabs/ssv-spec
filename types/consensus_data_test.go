package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestProposerConsensusDataJSONFieldsInSync guards the explicit ProposerConsensusData MarshalJSON/UnmarshalJSON
// overrides, which list the fields by hand to keep the JSON key order (Duty, Version, DataSSZ). A field added
// to the struct would otherwise silently vanish from JSON in both directions; if this fails, add the new
// field to both overrides in consensus_data.go too.
func TestProposerConsensusDataJSONFieldsInSync(t *testing.T) {
	require.Equal(t, 3, reflect.TypeOf(ProposerConsensusData{}).NumField(),
		"ProposerConsensusData gained a field — update its MarshalJSON/UnmarshalJSON overrides")
}
