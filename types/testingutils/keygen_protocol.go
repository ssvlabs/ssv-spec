package testingutils

import "github.com/bloxapp/ssv-spec/dkg"

type TestingKeygenProtocol struct {
	KeyGenOutput *dkg.KeyGenOutput
}

func (m TestingKeygenProtocol) Start(initOrReshare dkg.InitOrReshare) error {
	return nil
}

func (m TestingKeygenProtocol) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutcome, error) {
	return true, &dkg.KeyGenOutcome{KeyGenOutput: m.KeyGenOutput}, nil
}
