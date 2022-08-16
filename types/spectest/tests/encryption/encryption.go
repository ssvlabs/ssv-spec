package encryption

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/spectest/tests"
)

// SimpleEncrypt tests simple rsa encrypt
func SimpleEncrypt() *tests.EncryptionSpecTest {
	sk, pk, _ := types.GenerateKey()
	pkObj, _ := types.PemToPublicKey(pk)
	cipher, _ := types.Encrypt(pkObj, []byte("hello world"))
	return &tests.EncryptionSpecTest{
		Name:       "simple encryption",
		SKPem:      sk,
		PKPem:      pk,
		PlainText:  []byte("hello world"),
		CipherText: cipher,
	}
}
