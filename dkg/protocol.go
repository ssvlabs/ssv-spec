package dkg

import "github.com/bloxapp/ssv-spec/dkg/base"

// Protocol is an interface for all DKG protocol to support a variety of protocols for future upgrades
type Protocol interface {
	Start() ([]base.Message, error)
	// ProcessMsg returns true and a bls share if finished
	ProcessMsg(msg *base.Message) ([]base.Message, error)
	Output() ([]byte, error)
}
