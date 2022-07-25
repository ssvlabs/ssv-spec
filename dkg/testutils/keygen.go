package testutils

import dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"

type MockProtocol struct {
	LocalKeyShare *dkgtypes.LocalKeyShare
}

func (m MockProtocol) Start() ([]dkgtypes.Message, error) {
	return nil, nil
}

func (m MockProtocol) ProcessMsg(msg *dkgtypes.Message) ([]dkgtypes.Message, error) {
	return nil, nil
}

func (m MockProtocol) Output() ([]byte, error) {
	return m.LocalKeyShare.Encode()
}

func PlaceholderMessage() *dkgtypes.Message {
	return &dkgtypes.Message{
		Header: &dkgtypes.MessageHeader{
			SessionId: TestingRequestID[:],
			MsgType:   int32(dkgtypes.ProtocolMsgType),
			Sender:    uint64(1),
		},
	}
}
