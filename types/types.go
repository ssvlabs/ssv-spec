package types

import (
	"encoding/hex"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

type Root interface {
	GetRoot() ([32]byte, error)
}

type HexBLSPubKey spec.BLSPubKey
type HexBytes32 [32]byte
type HexBytes20 [20]byte
type HexBytes4 [4]byte

// MarshalJSON implements the json.Marshaler interface
func (h *HexBLSPubKey) MarshalJSON() ([]byte, error) {
	return marshalJson(h[:])
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (h *HexBLSPubKey) UnmarshalJSON(b []byte) error {
	return unmarshalJson(b, h[:])
}

// MarshalJSON implements the json.Marshaler interface
func (h *HexBytes32) MarshalJSON() ([]byte, error) {
	return marshalJson(h[:])
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (h *HexBytes32) UnmarshalJSON(b []byte) error {
	return unmarshalJson(b, h[:])
}

// MarshalJSON implements the json.Marshaler interface
func (h *HexBytes20) MarshalJSON() ([]byte, error) {
	return marshalJson(h[:])
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (h *HexBytes20) UnmarshalJSON(b []byte) error {
	return unmarshalJson(b, h[:])
}

// MarshalJSON implements the json.Marshaler interface
func (h *HexBytes4) MarshalJSON() ([]byte, error) {
	return marshalJson(h[:])
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (h *HexBytes4) UnmarshalJSON(b []byte) error {
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
