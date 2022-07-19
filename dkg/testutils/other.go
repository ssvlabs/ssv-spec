package testutils

import (
	"github.com/bloxapp/ssv-spec/dkg/base"
)

func PlaceholderMessage() *base.Message {
	return &base.Message{
		Header: &base.MessageHeader{
			MsgType: int32(base.ProtocolMsgType),
			Sender:  uint64(1),
		},
	}
}
