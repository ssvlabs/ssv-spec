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
