package testutils

import (
	"github.com/bloxapp/ssv-spec/dkg/types"
)

func PlaceholderMessage() *types.Message {
	return &types.Message{
		Header: &types.MessageHeader{
			MsgType: int32(types.ProtocolMsgType),
			Sender:  uint64(1),
		},
	}
}
