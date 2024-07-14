package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests a valid SignedSSVMessageTest
func Valid() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	// RSA key to sign message
	skByts, pkByts, err := types.GenerateKey()
	if err != nil {
		panic(err.Error())
	}
	sk, err := types.PemToPrivateKey(skByts)
	if err != nil {
		panic(err.Error())
	}

	msg := testingutils.TestingSignedSSVMessage(ks.Shares[1], 1, sk)

	return &SignedSSVMessageTest{
		Name: "valid",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		RSAPublicKey: [][]byte{pkByts},
	}
}
