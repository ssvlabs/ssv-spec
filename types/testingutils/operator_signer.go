package testingutils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/bloxapp/ssv-spec/types"
)

type testingOperatorSigner struct {
	SSVOperatorSK *rsa.PrivateKey
}

func NewTestingOperatorSigner(keySet *TestKeySet, operatorID types.OperatorID) *testingOperatorSigner {
	return &testingOperatorSigner{
		SSVOperatorSK: keySet.OperatorKeys[operatorID],
	}
}

func (km *testingOperatorSigner) SignSSVMessage(ssvMsg *types.SSVMessage) ([]byte, error) {
	encodedMsg, err := ssvMsg.Encode()
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(encodedMsg)
	signature, err := rsa.SignPKCS1v15(rand.Reader, km.SSVOperatorSK, crypto.SHA256, hash[:])
	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}
