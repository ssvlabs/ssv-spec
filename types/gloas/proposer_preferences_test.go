package gloas

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestProposerPreferences_JSON(t *testing.T) {
	p := &ProposerPreferences{
		DependentRoot:  phase0.Root{0x01, 0x02},
		ProposalSlot:   42,
		ValidatorIndex: 7,
		FeeRecipient:   bellatrix.ExecutionAddress{0xaa, 0xbb},
		TargetGasLimit: 36_000_000,
	}
	b, err := json.Marshal(p)
	require.NoError(t, err)
	// Lock the beacon-API wire form: snake_case keys, uint64 as decimal strings, root/address 0x-hex.
	require.JSONEq(t, fmt.Sprintf(
		`{"dependent_root":"%#x","proposal_slot":"42","validator_index":"7","fee_recipient":"%#x","target_gas_limit":"36000000"}`,
		p.DependentRoot, p.FeeRecipient), string(b))

	var dec ProposerPreferences
	require.NoError(t, json.Unmarshal(b, &dec))
	require.Equal(t, p, &dec)
}

func TestSignedProposerPreferences_JSON(t *testing.T) {
	s := &SignedProposerPreferences{
		Message: &ProposerPreferences{
			DependentRoot:  phase0.Root{0xaa},
			ProposalSlot:   9,
			ValidatorIndex: 3,
			FeeRecipient:   bellatrix.ExecutionAddress{0x11},
			TargetGasLimit: 30_000_000,
		},
		Signature: phase0.BLSSignature{0xbb, 0xcc},
	}
	b, err := json.Marshal(s)
	require.NoError(t, err)
	msgJSON, err := json.Marshal(s.Message)
	require.NoError(t, err)
	require.JSONEq(t, fmt.Sprintf(`{"message":%s,"signature":"%#x"}`, msgJSON, s.Signature), string(b))

	var dec SignedProposerPreferences
	require.NoError(t, json.Unmarshal(b, &dec))
	require.Equal(t, s, &dec)
}
