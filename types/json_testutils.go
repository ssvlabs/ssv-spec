package types

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

func (c *CommitteeID) UnmarshalJSON(data []byte) error {
	// Try base64-encoded string first
	var encoded string
	if err := json.Unmarshal(data, &encoded); err == nil {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return fmt.Errorf("failed to decode CommitteeID base64: %w", err)
		}
		if len(decoded) != 32 {
			return fmt.Errorf("invalid CommitteeID length from base64: expected 32, got %d", len(decoded))
		}
		copy(c[:], decoded)
		return nil
	}

	// If not a string, try array of integers
	var arr []byte
	if err := json.Unmarshal(data, &arr); err == nil {
		if len(arr) != 32 {
			return fmt.Errorf("invalid CommitteeID length from array: expected 32, got %d", len(arr))
		}
		copy(c[:], arr)
		return nil
	}

	// Fallback error
	return fmt.Errorf("CommitteeID must be base64 string or array of 32 bytes")
}

func (c *MessageID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		copy(c[:], str)
		return nil
	}

	// If not a string, try array of integers
	var arr []byte
	if err := json.Unmarshal(data, &arr); err == nil {
		if len(arr) != 56 {
			return fmt.Errorf("invalid MessageID length from array: expected 56, got %d", len(arr))
		}
		copy(c[:], arr)
		return nil
	}

	// Fallback error
	return fmt.Errorf("MessageID must be a string or array of 56 bytes")
}

func (d *DomainType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		copy(d[:], str)
		return nil
	}

	// If not a string, try array of integers
	var arr []byte
	if err := json.Unmarshal(data, &arr); err == nil {
		if len(arr) != 4 {
			return fmt.Errorf("invalid DomainType length from array: expected 4, got %d", len(arr))
		}
		copy(d[:], arr)
		return nil
	}

	// Fallback error
	return fmt.Errorf("DomainType must be a string or array of 4 bytes")
}

func (v *ValidatorPK) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		copy(v[:], str)
		return nil
	}

	var arr []byte
	if err := json.Unmarshal(data, &arr); err == nil {
		if len(arr) != 48 {
			return fmt.Errorf("invalid ValidatorPK length from array: expected 48, got %d", len(arr))
		}
		copy(v[:], arr)
		return nil
	}

	// Fallback error
	return fmt.Errorf("ValidatorPK must be a string or array of 48 bytes")
}

// UnmarshalJSON implements custom JSON unmarshaling for PartialSignatureMessage
func (p *PartialSignatureMessage) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to unmarshal into
	type tempPartialSignatureMessage struct {
		PartialSignature string `json:"PartialSignature"`
		SigningRoot      string `json:"SigningRoot"`
		Signer           uint64 `json:"Signer"`
		ValidatorIndex   string `json:"ValidatorIndex"`
	}

	var temp tempPartialSignatureMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal PartialSignatureMessage: %w", err)
	}

	// Convert PartialSignature from hex string to bytes
	partialSigHex := temp.PartialSignature
	if len(partialSigHex) > 2 && partialSigHex[:2] == "0x" {
		partialSigHex = partialSigHex[2:]
	}
	partialSigBytes, err := hex.DecodeString(partialSigHex)
	if err != nil {
		return fmt.Errorf("failed to decode PartialSignature hex: %w", err)
	}
	if len(partialSigBytes) != 96 {
		return fmt.Errorf("invalid PartialSignature length: expected 96, got %d", len(partialSigBytes))
	}
	p.PartialSignature = partialSigBytes

	// Convert SigningRoot from hex string to [32]byte
	signingRootHex := temp.SigningRoot
	if len(signingRootHex) > 2 && signingRootHex[:2] == "0x" {
		signingRootHex = signingRootHex[2:]
	}
	signingRootBytes, err := hex.DecodeString(signingRootHex)
	if err != nil {
		return fmt.Errorf("failed to decode SigningRoot hex: %w", err)
	}
	if len(signingRootBytes) != 32 {
		return fmt.Errorf("invalid SigningRoot length: expected 32, got %d", len(signingRootBytes))
	}
	copy(p.SigningRoot[:], signingRootBytes)

	// Set Signer
	p.Signer = OperatorID(temp.Signer)

	// Convert ValidatorIndex from string to phase0.ValidatorIndex
	validatorIndex, err := strconv.ParseUint(temp.ValidatorIndex, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse ValidatorIndex: %w", err)
	}
	p.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	return nil
}
