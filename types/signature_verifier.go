package types

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/pkg/errors"
)

// SignatureVerifier is an interface responsible for the verification of SignedSSVMessages
type SignatureVerifier interface {
	// Verify verifies a SignedSSVMessage's signature using the necessary keys extracted from the list of Operators
	Verify(msg *SignedSSVMessage, operators []*Operator) error
}

type RSAVerifier struct{}

func (r *RSAVerifier) Verify(msg *SignedSSVMessage, operators []*Operator) error {
	encodedMsg, err := msg.SSVMessage.Encode()
	if err != nil {
		return err
	}

	// Get message hash
	hash := sha256.Sum256(encodedMsg)

	// Find operator that matches ID with the signer and verify signature
	for i, signer := range msg.OperatorIDs {
		if err := r.VerifySignatureForSigner(hash, msg.Signatures[i], signer, operators); err != nil {
			return err
		}
	}
	return nil
}

func (r *RSAVerifier) VerifySignatureForSigner(root [32]byte, signature []byte, signer OperatorID,
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
