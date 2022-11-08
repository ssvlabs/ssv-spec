package frost

import (
	"bytes"

	ecies "github.com/ecies/go/v2"
)

func VerifyEciesKeyPair(pkbytes, skbytes []byte) bool {
	expectedplaintext := []byte("frost-ecies-test-keypair")
	sk := ecies.NewPrivateKeyFromBytes(skbytes)
	pk, _ := ecies.NewPublicKeyFromBytes(pkbytes)
	ciphertext, _ := ecies.Encrypt(pk, expectedplaintext)
	plaintext, _ := ecies.Decrypt(sk, ciphertext)
	return bytes.Equal(plaintext, expectedplaintext)
}
