package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests a valid SignedSSVMessageTest
func Valid() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	pkBytes, err := types.GetPublicKeyPem(ks.OperatorKeys[1])
	if err != nil {
		panic(err)
	}

	msg := testingutils.TestingSignedSSVMessage(ks.Shares[1], 1, ks.OperatorKeys[1])

	return &SignedSSVMessageTest{
		Name: "valid",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		RSAPublicKey: [][]byte{pkBytes},
	}
}
