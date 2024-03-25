package signedssvmsg

import (
	"testing"

	"crypto"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

type SignedSSVMessageTest struct {
	Name          string
	Messages      []*types.SignedSSVMessage
	ExpectedError string
	RSAPublicKey  [][]byte
}

func (test *SignedSSVMessageTest) TestName() string {
	return "signedssvmsg " + test.Name
}

func (test *SignedSSVMessageTest) Run(t *testing.T) {

	for _, msg := range test.Messages {

		// test validation
		err := msg.Validate()

		var data []byte
		if err == nil {
			// decode SSVMessage
			data, err = msg.SSVMessage.Encode()
		}

		// check RSA signature
		if err == nil {
			for i, rsaPublicKey := range test.RSAPublicKey {

				var pk *rsa.PublicKey
				pk, err = types.PemToPublicKey(rsaPublicKey)
				if err != nil {
					panic(err.Error())
				}

				messageHash := sha256.Sum256(data)
				err = rsa.VerifyPKCS1v15(pk, crypto.SHA256, messageHash[:], msg.Signature[i])
				if err != nil {
					break
				}
			}
		}

		if len(test.ExpectedError) != 0 {
			require.EqualError(t, err, test.ExpectedError)
		} else {
			require.NoError(t, err)
		}
	}
}