package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
