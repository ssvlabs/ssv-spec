package testingutils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/ssvlabs/ssv-spec/types"
)

type testingOperatorSigner struct {
	SSVOperatorSK *rsa.PrivateKey
	operatorID    types.OperatorID
}

func NewTestingOperatorSigner(keySet *TestKeySet, operatorID types.OperatorID) *testingOperatorSigner {
	return &testingOperatorSigner{
		SSVOperatorSK: keySet.OperatorKeys[operatorID],
		operatorID:    operatorID,
	}
}

func (km *testingOperatorSigner) SignSSVMessage(ssvMsg *types.SSVMessage) ([]byte, error) {
	return SignSSVMessage(km.SSVOperatorSK, ssvMsg)
}

// GetOperatorID returns the operator ID
func (km *testingOperatorSigner) GetOperatorID() types.OperatorID {
	return km.operatorID
}

func SignSSVMessage(sk *rsa.PrivateKey, ssvMsg *types.SSVMessage) ([]byte, error) {
	encodedMsg, err := ssvMsg.Encode()
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(encodedMsg)
	signature, err := rsa.SignPKCS1v15(rand.Reader, sk, crypto.SHA256, hash[:])
	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}
