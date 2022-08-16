package encryption

import (
	"fmt"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/spectest/tests"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// EncryptBLSSK tests encrypting a BLS private key
func EncryptBLSSK() *tests.EncryptionSpecTest {
	types.InitBLS()
	blsSK := &bls.SecretKey{}
	blsSK.SetByCSPRNG()

	sk, pk, _ := types.GenerateKey()
	pkObj, _ := types.PemToPublicKey(pk)
	cipher, _ := types.Encrypt(pkObj, blsSK.Serialize())

	fmt.Printf("cipher L: %d\n", len(cipher))
	return &tests.EncryptionSpecTest{
		Name:       "bls secret key encryption",
		SKPem:      sk,
		PKPem:      pk,
		PlainText:  blsSK.Serialize(),
		CipherText: cipher,
	}
}
