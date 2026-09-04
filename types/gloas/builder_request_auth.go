package gloas

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// BuilderRequestAuth authenticates a cluster's direct request to a builder for a proposal slot (SIP #94
// §5). A validator-key BLS signature over it under DomainBuilderRequestAuth
// — chain-independent (genesis compute_domain, like the deprecated ValidatorRegistrationV1) — yields
// SignedBuilderRequestAuth. Data is the token agreed with the builder out of band (defaulting to the
// builder's advertised URL bytes; zero-length is invalid), signed byte-for-byte. Slot is the proposal
// slot the request is authorized for, never the slot it is signed at — a pure function of the proposer
// lookahead, carrying no dependent root, so re-emissions reproduce a byte-identical signing root.
type BuilderRequestAuth struct {
	Data []byte `ssz-max:"4096"`
	Slot phase0.Slot
}

// SignedBuilderRequestAuth is a BuilderRequestAuth plus the reconstructed validator signature, held for
// the proposal slot and forwarded to the builder as each entry's auth in the §4 produce body (and,
// optionally, via the builder_preferences endpoint).
type SignedBuilderRequestAuth struct {
	Message   *BuilderRequestAuth
	Signature phase0.BLSSignature `ssz-size:"96"`
}

// builderRequestAuthJSON is the beacon-API JSON form: data as 0x-hex, slot as a decimal string.
type builderRequestAuthJSON struct {
	Data string `json:"data"`
	Slot string `json:"slot"`
}

// MarshalJSON implements json.Marshaler.
func (b *BuilderRequestAuth) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderRequestAuthJSON{
		Data: fmt.Sprintf("%#x", b.Data),
		Slot: fmt.Sprintf("%d", b.Slot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BuilderRequestAuth) UnmarshalJSON(input []byte) error {
	var data builderRequestAuthJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	raw, err := hex.DecodeString(strings.TrimPrefix(data.Data, "0x"))
	if err != nil {
		return fmt.Errorf("invalid value for data: %w", err)
	}
	b.Data = raw
	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid value for slot: %w", err)
	}
	b.Slot = phase0.Slot(slot)
	return nil
}

// signedBuilderRequestAuthJSON is the beacon-API JSON form of SignedBuilderRequestAuth.
type signedBuilderRequestAuthJSON struct {
	Message   *BuilderRequestAuth `json:"message"`
	Signature string              `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedBuilderRequestAuth) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBuilderRequestAuthJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedBuilderRequestAuth) UnmarshalJSON(input []byte) error {
	var data signedBuilderRequestAuthJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if data.Message == nil {
		return errors.New("message missing")
	}
	s.Message = data.Message
	if err := decodeHexInto(s.Signature[:], data.Signature, "signature"); err != nil {
		return err
	}
	return nil
}
