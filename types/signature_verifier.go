package types

// SignatureVerifier is an interface responsible for the verification of SignedSSVMessages
type SignatureVerifier interface {
	// Verify verifies a SignedSSVMessage's signature using the necessary keys extracted from the list of Operators
	Verify(msg *SignedSSVMessage, operators []*Operator) error
}

type RSAVerifier struct{}
