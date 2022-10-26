package testingutils

import "github.com/bloxapp/ssv-spec/dkg"

type TestingKeygenProtocol struct {
	KeyGenOutput *dkg.KeyGenOutput
}

func (m TestingKeygenProtocol) Start() error {
	return nil
}

func (m TestingKeygenProtocol) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {
	return true, &dkg.ProtocolOutcome{ProtocolOutput: m.KeyGenOutput}, nil
}
