package types

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type OperatorSigner struct {
	SSVOperatorSK *rsa.PrivateKey
	OperatorID    OperatorID
}

func (km *OperatorSigner) SignSSVMessage(ssvMsg *SSVMessage) ([]byte, error) {
	return SignSSVMessage(km.SSVOperatorSK, ssvMsg)
}

// GetOperatorID returns the operator ID
func (km *OperatorSigner) GetOperatorID() OperatorID {
	return km.OperatorID
}

func SignSSVMessage(sk *rsa.PrivateKey, ssvMsg *SSVMessage) ([]byte, error) {
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
