package drand

import (
	"encoding/json"
	"github.com/drand/kyber/share/dkg"
)

type MsgType int

const (
	DealBundleMsg MsgType = iota
	ResponseBundleMsg
	JustificationBundleMsg
)

type Message struct {
	MsgType
	DealBundle          *dkg.DealBundle
	ResponseBundle      *dkg.ResponseBundle
	JustificationBundle *dkg.JustificationBundle
}

func (msg *Message) validate() bool {
	return true
}

// Encode returns a msg encoded bytes or error
func (msg *Message) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *Message) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}
