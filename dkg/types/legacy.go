package types

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/pkg/errors"
)

type LegacyMessage struct {
	MsgType    MsgType
	Identifier RequestID
	Data       []byte
}

// Encode returns a msg encoded bytes or error
func (msg *LegacyMessage) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *LegacyMessage) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

func (msg *LegacyMessage) Validate() error {
	// TODO msg type
	// TODO len(data)
	return nil
}

func (msg *LegacyMessage) GetRoot() ([]byte, error) {
	marshaledRoot, err := msg.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode PartialSignatureMessage")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}
