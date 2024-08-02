package p2p

import "github.com/ssvlabs/ssv-spec/types"

// Broadcaster is the interface used to abstract message broadcasting
type Broadcaster interface {
	Broadcast(msgID types.MessageID, message *types.SignedSSVMessage) error
	BroadcastRaw([]string, []byte) error
}

// Subscriber is used to abstract topic management
type Subscriber interface {
	Subscribe(vpk types.ValidatorPK) error
}
