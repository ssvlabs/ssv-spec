package qbft

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

// This file adds, as testing utils, the Encode, Decode and GetRoot methods
// so that structures follow the types.Encoder and types.Root interface

// Controller
func (c *Controller) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Controller) Decode(data []byte) error {
	err := json.Unmarshal(data, c)
	if err != nil {
		return errors.Wrap(err, "could not decode controller")
	}

	config := c.GetConfig()
	for _, i := range c.StoredInstances {
		if i != nil {
			i.config = config
		}
	}
	return nil
}

func (c *Controller) GetRoot() ([32]byte, error) {
	marshaledRoot, err := json.Marshal(c)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode controller")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// UnmarshalJSON is a custom JSON unmarshaller for Controller
func (c *Controller) UnmarshalJSON(data []byte) error {
	type ControllerAlias Controller
	aux := &struct {
		*ControllerAlias
	}{
		ControllerAlias: (*ControllerAlias)(c),
	}

	return json.Unmarshal(data, aux)
}

// Instance
func (i *Instance) Encode() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Instance) Decode(data []byte) error {
	return json.Unmarshal(data, &i)
}

func (i *Instance) GetRoot() ([32]byte, error) {
	return i.State.GetRoot()
}

// MarshalJSON is a custom JSON marshaller for Instance
func (i *Instance) MarshalJSON() ([]byte, error) {
	type Alias Instance
	if i.forceStop {
		return json.Marshal(&struct {
			ForceStop bool `json:"forceStop"`
			*Alias
		}{
			ForceStop: i.forceStop,
			Alias:     (*Alias)(i),
		})
	} else {
		return json.Marshal(&struct {
			*Alias
		}{
			Alias: (*Alias)(i),
		})
	}
}

// UnmarshalJSON is a custom JSON unmarshaller for Instance
func (i *Instance) UnmarshalJSON(data []byte) error {
	type Alias Instance
	aux := &struct {
		ForceStop *bool `json:"forceStop,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.ForceStop != nil {
		i.forceStop = *aux.ForceStop
	}
	return nil
}

// State

func (s *State) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *State) Decode(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s *State) GetRoot() ([32]byte, error) {
	marshaledRoot, err := s.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode state")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// MsgContainer
func (c *MsgContainer) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func (c *MsgContainer) Decode(data []byte) error {
	return json.Unmarshal(data, &c)
}

func (c *MsgContainer) GetRoot() ([32]byte, error) {
	marshaledRoot, err := c.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode msg container")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// Message
func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"MsgType":                  m.MsgType,
		"Height":                   m.Height,
		"Round":                    m.Round,
		"Identifier":               m.Identifier,
		"Root":                     hex.EncodeToString(m.Root[:]),
		"DataRound":                m.DataRound,
		"RoundChangeJustification": m.RoundChangeJustification,
		"PrepareJustification":     m.PrepareJustification,
	})
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		Root string `json:"Root"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Root != "" {
		if bytes, err := hex.DecodeString(aux.Root); err == nil {
			if len(bytes) != 32 {
				return errors.New("Root must be exactly 32 bytes")
			}
			copy(m.Root[:], bytes)
			return nil
		}
	}

	return nil
}

// Value type for hex encoding json

type Value [32]byte

func (r *Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(r[:]))
}

func (r *Value) UnmarshalJSON(data []byte) error {
	hexStr := hexStringFromJSON(data)
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return errors.Wrap(err, "failed to decode Value")
	}
	copy(r[:], bytes)
	return nil
}

// helpers

// hexStringFromJSON trims surrounding quotes from a JSON string value.
func hexStringFromJSON(data []byte) string {
	return strings.TrimSuffix(strings.TrimPrefix(string(data), "\""), "\"")
}
