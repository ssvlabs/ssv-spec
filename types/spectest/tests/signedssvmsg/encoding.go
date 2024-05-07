package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Encoding tests encoding of a SignedSSVMessage
func Encoding() *EncodingTest {

	ks := testingutils.Testing4SharesSet()

	// RSA key to sign message
	skByts, _, err := types.GenerateKey()
	if err != nil {
		panic(err.Error())
	}
	sk, err := types.PemToPrivateKey(skByts)
	if err != nil {
		panic(err.Error())
	}

	msg := testingutils.TestingSignedSSVMessage(ks.Shares[1], 1, sk)

	byts, err := msg.Encode()
	if err != nil {
		panic(err.Error())
	}

	return &EncodingTest{
		Name: "encoding",
		Data: byts,
	}
}
