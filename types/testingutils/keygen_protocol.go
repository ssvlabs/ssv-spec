package testingutils

import "github.com/bloxapp/ssv-spec/dkg"

type MockKeygenProtocol struct {
	KeyGenOutput *dkg.KeyGenOutput
}

func (m MockKeygenProtocol) Start(init *dkg.Init) error {
	return nil
}

func (m MockKeygenProtocol) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error) {
	return true, m.KeyGenOutput, nil
}
