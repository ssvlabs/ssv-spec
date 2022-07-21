package testutils

import dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"

type mockProtocol struct {
	localKeyShare *dkgtypes.LocalKeyShare
}

func (m mockProtocol) Start() ([]dkgtypes.Message, error) {
	return nil, nil
}

func (m mockProtocol) ProcessMsg(msg *dkgtypes.Message) ([]dkgtypes.Message, error) {
	return nil, nil
}

func (m mockProtocol) Output() ([]byte, error) {
	return m.localKeyShare.Encode()
}
