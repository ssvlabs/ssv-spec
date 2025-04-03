package signedssvmsg

import (
	"testing"

	"crypto"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/ssvlabs/ssv-spec/types"
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

		// encode message
		var encodedMsg []byte
		if err == nil {
			encodedMsg, err = msg.SSVMessage.Encode()
		}

		// check RSA signature
		if err == nil {

			messageHash := sha256.Sum256(encodedMsg)

			for i, pkBytes := range test.RSAPublicKey {
				var pk *rsa.PublicKey
				pk, err = types.PemToPublicKey(pkBytes)
				if err != nil {
					panic(err.Error())
				}

				err = rsa.VerifyPKCS1v15(pk, crypto.SHA256, messageHash[:], msg.Signatures[i])
			}
		}

		if len(test.ExpectedError) != 0 {
			require.EqualError(t, err, test.ExpectedError)
		} else {
			require.NoError(t, err)
		}
	}
}
