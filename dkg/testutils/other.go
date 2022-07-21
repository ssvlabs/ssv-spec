package testutils

import (
	"encoding/hex"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
)

var TestingWithdrawalCredentials, _ = hex.DecodeString("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f")
var TestingForkVersion = types.PraterNetwork.ForkVersion()

func PlaceholderMessage() *dkgtypes.Message {
	return &dkgtypes.Message{
		Header: &dkgtypes.MessageHeader{
			MsgType: int32(dkgtypes.ProtocolMsgType),
			Sender:  uint64(1),
		},
	}
}

func FakeEncryption(data []byte) []byte {
	out := []byte("__fake_encrypted(")
	out = append(out, data...)
	out = append(out, []byte(")")...)
	return out
}
