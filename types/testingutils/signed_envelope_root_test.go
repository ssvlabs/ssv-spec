package testingutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	eth2gloas "github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// TestSignedEnvelopeRootMatchesGoEth2Client pins what the testing beacon node records for a published §6 reveal
// (SubmitExecutionPayloadEnvelope): the root of go-eth2-client's SignedExecutionPayloadEnvelope holding the full
// envelope the blinded one stands for. The vectors expect that root, and an implementation compares it with the
// root of the signed envelope its own beacon node mock receives.
func TestSignedEnvelopeRootMatchesGoEth2Client(t *testing.T) {
	hex := func(seed byte, n int) string { return fmt.Sprintf("%#x", bytes.Repeat([]byte{seed}, n)) }
	// The payload is built through its JSON form: its base fee is a type from a module this one does not import
	// directly.
	payloadJSON, err := json.Marshal(map[string]any{
		"parent_hash":       hex(0xa0, 32),
		"fee_recipient":     hex(0xa1, 20),
		"state_root":        hex(0xa2, 32),
		"receipts_root":     hex(0xa3, 32),
		"logs_bloom":        hex(0xa4, 256),
		"prev_randao":       hex(0xa5, 32),
		"block_number":      "100",
		"gas_limit":         "36000000",
		"gas_used":          "21000",
		"timestamp":         "1700000000",
		"extra_data":        hex(0xa6, 8),
		"base_fee_per_gas":  "7",
		"block_hash":        hex(0xa7, 32),
		"transactions":      []string{hex(0xa8, 40), hex(0xa9, 3)},
		"withdrawals":       []map[string]string{{"index": "1", "validator_index": "2", "address": hex(0xab, 20), "amount": "3"}},
		"blob_gas_used":     "131072",
		"excess_blob_gas":   "262144",
		"block_access_list": hex(0xad, 50),
		"slot_number":       "66",
	})
	require.NoError(t, err)
	payload := &eth2gloas.ExecutionPayload{}
	require.NoError(t, json.Unmarshal(payloadJSON, payload))

	full := &eth2gloas.ExecutionPayloadEnvelope{
		Payload:               payload,
		ExecutionRequests:     &eth2gloas.ExecutionRequests{},
		BuilderIndex:          eth2gloas.BuilderIndex(gloas.BuilderIndexSelfBuild),
		BeaconBlockRoot:       phase0.Root{0xb0},
		ParentBeaconBlockRoot: phase0.Root{0xb1},
	}
	signature := phase0.BLSSignature{0xc0, 0xc1, 0xc2}
	want, err := (&eth2gloas.SignedExecutionPayloadEnvelope{Message: full, Signature: signature}).HashTreeRoot()
	require.NoError(t, err)

	payloadRoot, err := full.Payload.HashTreeRoot()
	require.NoError(t, err)
	requestsRoot, err := full.ExecutionRequests.HashTreeRoot()
	require.NoError(t, err)
	blinded := &gloas.BlindedExecutionPayloadEnvelope{
		PayloadRoot:           payloadRoot,
		ExecutionRequestsRoot: requestsRoot,
		BuilderIndex:          gloas.BuilderIndex(full.BuilderIndex),
		BeaconBlockRoot:       full.BeaconBlockRoot,
		ParentBeaconBlockRoot: full.ParentBeaconBlockRoot,
	}
	got, err := signedEnvelopeRoot(blinded, signature)
	require.NoError(t, err)

	require.Equal(t, phase0.Root(want), got)
}
