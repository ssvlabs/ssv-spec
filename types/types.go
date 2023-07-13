package types

import (
	"encoding/hex"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

type Root interface {
	GetRoot() ([32]byte, error)
}

type HexBLSPubKey spec.BLSPubKey
type HexBytes32 [32]byte
type HexExecutionAddress bellatrix.ExecutionAddress

// MarshalJSON implements the json.Marshaler interface
func (h HexBLSPubKey) MarshalJSON() ([]byte, error) {
	return marshalJson(h[:])
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (h HexBLSPubKey) UnmarshalJSON(b []byte) error {
	return unmarshalJson(b, h[:])
}

// MarshalJSON implements the json.Marshaler interface
func (h HexBytes32) MarshalJSON() ([]byte, error) {
	return marshalJson(h[:])
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (h HexBytes32) UnmarshalJSON(b []byte) error {
	return unmarshalJson(b, h[:])
}

// MarshalJSON implements the json.Marshaler interface
func (h *HexExecutionAddress) MarshalJSON() ([]byte, error) {
	return marshalJson(h[:])
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (h *HexExecutionAddress) UnmarshalJSON(b []byte) error {
	return unmarshalJson(b, h[:])
}

func marshalJson(h []byte) ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(h)))
	hex.Encode(dst, h[:])
	return json.Marshal(string(dst))
}

func unmarshalJson(b []byte, h []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	decoded, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	copy(h[:], decoded)
	return nil
}
