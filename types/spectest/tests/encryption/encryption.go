package encryption

import (
	"github.com/bloxapp/ssv-spec/types"
)

// SimpleEncrypt tests simple rsa encrypt
func SimpleEncrypt() *EncryptionSpecTest {
	sk, pk, _ := types.GenerateKey()
	pkObj, _ := types.PemToPublicKey(pk)
	cipher, _ := types.Encrypt(pkObj, []byte("hello world"))
	return &EncryptionSpecTest{
		Name:       "simple encryption",
		SKPem:      sk,
		PKPem:      pk,
		PlainText:  []byte("hello world"),
		CipherText: cipher,
	}
}
