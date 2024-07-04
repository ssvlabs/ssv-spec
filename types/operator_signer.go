package types

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/pkg/errors"
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

func Verify(msg *SignedSSVMessage, operators []*Operator) error {
	encodedMsg, err := msg.SSVMessage.Encode()
	if err != nil {
		return err
	}

	// Get message hash
	hash := sha256.Sum256(encodedMsg)

	// Find operator that matches ID with the signer and verify signature
	for i, signer := range msg.OperatorIDs {
		if err := verifySignatureForSigner(hash, msg.Signatures[i], signer, operators); err != nil {
			return err
		}
	}
	return nil
}

func verifySignatureForSigner(root [32]byte, signature []byte, signer OperatorID,
	operators []*Operator) error {
	for _, op := range operators {
		// Find signer
		if signer == op.OperatorID {
			// Get public key
			pk, err := PemToPublicKey(op.SSVOperatorPubKey)
			if err != nil {
				return errors.Wrap(err, "could not parse signer public key")
			}

			// Verify
			err = rsa.VerifyPKCS1v15(pk, crypto.SHA256, root[:], signature)

			return err
		}
	}
	return errors.New("unknown signer")
}
