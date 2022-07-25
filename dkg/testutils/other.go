package testutils

import (
	"encoding/hex"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
)

var TestingWithdrawalCredentials, _ = hex.DecodeString("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f")
var TestingForkVersion = types.PraterNetwork.ForkVersion()
var TestingAddress = common.HexToAddress("535953b5a6040074948cf185eaa7d2abbd66808f")
var TestingRequestID = dkgtypes.NewRequestID(TestingAddress, 1)

func PlaceholderMessage() *dkgtypes.Message {
	return &dkgtypes.Message{
		Header: &dkgtypes.MessageHeader{
			SessionId: TestingRequestID[:],
			MsgType:   int32(dkgtypes.ProtocolMsgType),
			Sender:    uint64(1),
		},
	}
}

func FakeEncryption(data []byte) []byte {
	out := []byte("__fake_encrypted(")
	out = append(out, data...)
	out = append(out, []byte(")")...)
	return out
}

func FakeEcdsaSign(root []byte, address []byte) []byte {
	out := []byte("__fake_ecdsa_sign(root=")
	out = append(out, root...)
	out = append(out, []byte(",address=")...)
	out = append(out, address...)
	out = append(out, []byte(")")...)
	return out
}
