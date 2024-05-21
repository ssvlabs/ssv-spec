package encryption

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/ssvlabs/ssv-spec/types"
)

// EncryptBLSSK tests encrypting a BLS private key
func EncryptBLSSK() *EncryptionSpecTest {
	types.InitBLS()
	blsSK := &bls.SecretKey{}
	blsSK.SetByCSPRNG()

	sk, pk, _ := types.GenerateKey()
	pkObj, _ := types.PemToPublicKey(pk)
	cipher, _ := types.Encrypt(pkObj, blsSK.Serialize())

	fmt.Printf("cipher L: %d\n", len(cipher))
	return &EncryptionSpecTest{
		Name:       "bls secret key encryption",
		SKPem:      sk,
		PKPem:      pk,
		PlainText:  blsSK.Serialize(),
		CipherText: cipher,
	}
}
