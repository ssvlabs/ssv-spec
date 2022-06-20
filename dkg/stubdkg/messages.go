package stubdkg

import (
	"encoding/json"
)

type Stage int

const (
	StubStage1 Stage = iota
	StubStage2
	StubStage3
)

type ProtocolMsg struct {
	Stage Stage
	// Data is any data a real DKG implementation will need
	Data interface{}
}

// Encode returns a msg encoded bytes or error
func (msg *ProtocolMsg) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *ProtocolMsg) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}
