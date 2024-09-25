package encryption

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SimpleEncrypt tests simple rsa encrypt
func SimpleEncrypt() *EncryptionSpecTest {

	ks := testingutils.Testing4SharesSet()

	sk := ks.OperatorKeys[1]
	skPem := types.PrivateKeyToPem(sk)
	pkPem, err := types.GetPublicKeyPem(sk)
	if err != nil {
		panic(err)
	}

	return &EncryptionSpecTest{
		Name:      "simple encryption",
		SKPem:     skPem,
		PKPem:     pkPem,
		PlainText: []byte("hello world"),
	}
}
