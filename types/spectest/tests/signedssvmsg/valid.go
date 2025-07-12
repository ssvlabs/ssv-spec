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

	return NewSignedSSVMessageTest(
		"valid",
		"Test validation of a valid signed SSV message with proper RSA signature",
		[]*types.SignedSSVMessage{msg},
		"",
		[][]byte{pkBytes},
	)
}
