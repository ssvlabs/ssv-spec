package encryption

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

type EncryptionSpecTest struct {
	Name      string
	SKPem     []byte
	PKPem     []byte
	PlainText []byte
}

func (test *EncryptionSpecTest) TestName() string {
	return "encryption " + test.Name
}

func (test *EncryptionSpecTest) Run(t *testing.T) {
	// get sk from pem
	sk, err := types.PemToPrivateKey(test.SKPem)
	require.NoError(t, err)

	// get pk from sk and compare to test pk
	pkFromSK, err := types.GetPublicKeyPem(sk)
	require.NoError(t, err)
	require.EqualValues(t, test.PKPem, pkFromSK)

	pk, err := types.PemToPublicKey(test.PKPem)
	require.NoError(t, err)

	// encrypt
	cipher, err := types.Encrypt(pk, test.PlainText)
	require.NoError(t, err)

	// decrypt and compare to plain text
	plain, err := types.Decrypt(sk, cipher)
	require.NoError(t, err)
	require.EqualValues(t, test.PlainText, plain)
}
