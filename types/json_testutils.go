package types

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ************** ExpectedRoot **************
type ExpectedRoot [32]byte

func (r *ExpectedRoot) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(r[:]))
}

func (r *ExpectedRoot) UnmarshalJSON(data []byte) error {
	hexStr := hexStringFromJSON(data)
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return fmt.Errorf("failed to decode ExpectedRoot: %w", err)
	}
	copy(r[:], bytes)
	return nil
}

// ************** ExpectedCdRoot **************
type ExpectedCdRoot [32]byte

func (r *ExpectedCdRoot) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(r[:]))
}

func (r *ExpectedCdRoot) UnmarshalJSON(data []byte) error {
	hexStr := hexStringFromJSON(data)
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return fmt.Errorf("failed to decode ExpectedCdRoot: %w", err)
	}
	copy(r[:], bytes)
	return nil
}

// ************** ExpectedBlkRoot **************
type ExpectedBlkRoot [32]byte

func (r *ExpectedBlkRoot) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(r[:]))
}

func (r *ExpectedBlkRoot) UnmarshalJSON(data []byte) error {
	hexStr := hexStringFromJSON(data)
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return fmt.Errorf("failed to decode ExpectedBlkRoot: %w", err)
	}
	copy(r[:], bytes)
	return nil
}

// ************** ForkVersion **************

type ForkVersion [4]byte

func (f *ForkVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(f[:]))
}

func (f *ForkVersion) UnmarshalJSON(data []byte) error {
	hexStr := hexStringFromJSON(data)
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return fmt.Errorf("failed to decode ForkVersion: %w", err)
	}
	copy(f[:], bytes)
	return nil
}

// ************** CommitteeID **************

func (c *CommitteeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(c[:]))
}

func (c *CommitteeID) UnmarshalJSON(data []byte) error {
	hexStr := hexStringFromJSON(data)
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		bytes, err = tryUnmarshalByteArray(data, 32)
		if err != nil {
			return fmt.Errorf("failed to decode CommitteeID: %w", err)
		}
	}
	copy(c[:], bytes)
	return nil
}

// ************** MessageID **************

func (c *MessageID) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(c[:]))
}

func (c *MessageID) UnmarshalJSON(data []byte) error {
	hexStr := hexStringFromJSON(data)
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		bytes, err = tryUnmarshalByteArray(data, 56)
		if err != nil {
			return fmt.Errorf("failed to decode MessageID: %w", err)
		}
	}
	copy(c[:], bytes)
	return nil
}

// ************** DomainType **************

func (d *DomainType) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(d[:]))
}

func (d *DomainType) UnmarshalJSON(data []byte) error {
	hexStr := hexStringFromJSON(data)
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		bytes, err = tryUnmarshalByteArray(data, 4)
		if err != nil {
			return fmt.Errorf("failed to decode DomainType: %w", err)
		}
	}
	copy(d[:], bytes)
	return nil
}

// ************** PartialSignatureMessage **************

func (p *PartialSignatureMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"PartialSignature": hex.EncodeToString(p.PartialSignature),
		"SigningRoot":      hex.EncodeToString(p.SigningRoot[:]),
		"Signer":           p.Signer,
		"ValidatorIndex":   strconv.FormatUint(uint64(p.ValidatorIndex), 10),
	})
}

func (p *PartialSignatureMessage) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to unmarshal into
	type tempPartialSignatureMessage struct {
		PartialSignature string      `json:"PartialSignature"`
		SigningRoot      interface{} `json:"SigningRoot"`
		Signer           uint64      `json:"Signer"`
		ValidatorIndex   string      `json:"ValidatorIndex"`
	}

	var temp tempPartialSignatureMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal PartialSignatureMessage: %w", err)
	}

	if temp.PartialSignature == "" {
		p.PartialSignature = nil
	} else {
		partialSigHex := temp.PartialSignature
		partialSigHex = strings.TrimPrefix(partialSigHex, "0x")
		partialSigBytes, err := hex.DecodeString(partialSigHex)
		if err != nil {
			partialSigBytes, err = base64.StdEncoding.DecodeString(temp.PartialSignature)
			if err != nil {
				return fmt.Errorf("failed to decode PartialSignature (tried hex and base64): %w", err)
			}
		}
		p.PartialSignature = partialSigBytes
	}

	// Convert SigningRoot from hex string or byte array to [32]byte
	var signingRootBytes []byte
	switch signingRoot := temp.SigningRoot.(type) {
	case string:
		// Handle hex string
		signingRootHex := signingRoot
		if len(signingRootHex) > 2 && signingRootHex[:2] == "0x" {
			signingRootHex = signingRootHex[2:]
		}
		var err error
		signingRootBytes, err = hex.DecodeString(signingRootHex)
		if err != nil {
			return fmt.Errorf("failed to decode SigningRoot hex: %w", err)
		}
	case []interface{}:
		// Handle byte array
		signingRootBytes = make([]byte, len(signingRoot))
		for i, b := range signingRoot {
			switch v := b.(type) {
			case float64:
				signingRootBytes[i] = byte(v)
			case int:
				signingRootBytes[i] = byte(v)
			default:
				return fmt.Errorf("invalid byte type in SigningRoot array: %T", b)
			}
		}
	case []byte:
		// Handle direct byte array
		signingRootBytes = signingRoot
	default:
		return fmt.Errorf("invalid SigningRoot type: %T", temp.SigningRoot)
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

// ************** Share **************

func (s *Share) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ValidatorIndex":      s.ValidatorIndex,
		"ValidatorPubKey":     hex.EncodeToString(s.ValidatorPubKey[:]),
		"SharePubKey":         s.SharePubKey,
		"Committee":           s.Committee,
		"DomainType":          s.DomainType,
		"FeeRecipientAddress": hex.EncodeToString(s.FeeRecipientAddress[:]),
		"Graffiti":            s.Graffiti,
	})
}

func (s *Share) UnmarshalJSON(data []byte) error {
	type tempShare struct {
		ValidatorIndex      phase0.ValidatorIndex
		ValidatorPubKey     string           `json:"ValidatorPubKey"`
		SharePubKey         ShareValidatorPK `ssz-size:"48"`
		Committee           []*ShareMember   `ssz-max:"13"`
		DomainType          DomainType       `ssz-size:"4"`
		FeeRecipientAddress string           `json:"FeeRecipientAddress"`
		Graffiti            []byte           `json:"Graffiti"`
	}

	var temp tempShare
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal Share: %w", err)
	}

	feeRecipientAddress, err := hex.DecodeString(temp.FeeRecipientAddress)
	if err != nil {
		return fmt.Errorf("failed to decode FeeRecipientAddress: %w", err)
	}

	validatorPubKey, err := hex.DecodeString(temp.ValidatorPubKey)
	if err != nil {
		return fmt.Errorf("failed to decode ValidatorPubKey: %w", err)
	}

	s.ValidatorIndex = temp.ValidatorIndex
	copy(s.ValidatorPubKey[:], validatorPubKey)
	s.SharePubKey = temp.SharePubKey
	s.Committee = temp.Committee
	s.DomainType = temp.DomainType
	copy(s.FeeRecipientAddress[:], feeRecipientAddress)
	s.Graffiti = temp.Graffiti

	return nil
}

// hexStringFromJSON trims surrounding quotes from a JSON string value.
func hexStringFromJSON(data []byte) string {
	return strings.TrimSuffix(strings.TrimPrefix(string(data), "\""), "\"")
}

// tryUnmarshalByteArray tries to unmarshal a JSON array of bytes of a given length.
func tryUnmarshalByteArray(data []byte, expectedLen int) ([]byte, error) {
	var arr []byte
	if err := json.Unmarshal(data, &arr); err == nil {
		if len(arr) != expectedLen {
			return nil, fmt.Errorf("invalid array length: expected %d, got %d", expectedLen, len(arr))
		}
		return arr, nil
	}
	return nil, fmt.Errorf("not a valid byte array")
}
