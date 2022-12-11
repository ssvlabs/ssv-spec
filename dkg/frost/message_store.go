package frost

import "github.com/bloxapp/ssv-spec/dkg"

type ProtocolMessageStore map[ProtocolRound]map[uint32]*dkg.SignedMessage

func newProtocolMessageStore() ProtocolMessageStore {
	m := make(map[ProtocolRound]map[uint32]*dkg.SignedMessage)
	for _, round := range rounds {
		m[round] = make(map[uint32]*dkg.SignedMessage)
	}
	return m
}

func (msgStore ProtocolMessageStore) allMessagesReceivedFor(round ProtocolRound, operators []uint32) bool {
	for _, operatorID := range operators {
		if _, ok := msgStore[round][operatorID]; !ok {
			return false
		}
	}
	return true
}
