package tests

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NetworkBehaviour configures non-default testing network responses. Like BeaconNodeBehaviour it is part of the
// vector JSON, so an implementation running the vector configures its own network mock the same way.
type NetworkBehaviour struct {
	// FailPartialSigBroadcast makes the broadcast of a partial-signature container of any of these types fail
	// (SIP #94 §5).
	FailPartialSigBroadcast []types.PartialSigMsgType `json:",omitempty"`
}

// Apply configures net; a nil behaviour leaves the defaults.
func (b *NetworkBehaviour) Apply(net *testingutils.TestingNetwork) {
	if b == nil {
		return
	}
	net.SetFailPartialSigBroadcast(b.FailPartialSigBroadcast)
}
