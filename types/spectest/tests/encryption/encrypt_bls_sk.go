package encryption

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// EncryptBLSSK tests encrypting a BLS private key
func EncryptBLSSK() *EncryptionSpecTest {

	ks := testingutils.Testing4SharesSet()

	sk := ks.OperatorKeys[1]
	skPem := types.PrivateKeyToPem(sk)
	pkPem, err := types.GetPublicKeyPem(sk)
	if err != nil {
		panic(err)
	}

	blsSK := ks.Shares[1]

	return NewEncryptionSpecTest(
		"bls secret key encryption",
		skPem,
		pkPem,
		blsSK.Serialize(),
	)
}
