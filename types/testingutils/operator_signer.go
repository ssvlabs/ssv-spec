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
}

func NewTestingOperatorSigner(keySet *TestKeySet, operatorID types.OperatorID) *testingOperatorSigner {
	return &testingOperatorSigner{
		SSVOperatorSK: keySet.OperatorKeys[operatorID],
	}
}

func (km *testingOperatorSigner) SignSSVMessage(data []byte) ([256]byte, error) {
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, km.SSVOperatorSK, crypto.SHA256, hash[:])
	if err != nil {
		return [256]byte{}, err
	}

	sig := [256]byte{}
	copy(sig[:], signature)

	return sig, nil
}
