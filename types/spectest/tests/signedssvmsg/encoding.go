package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Encoding tests encoding of a SignedSSVMessage
func Encoding() *EncodingTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingSignedSSVMessage(ks.Shares[1], 1, ks.OperatorKeys[1])

	byts, err := msg.Encode()
	if err != nil {
		panic(err.Error())
	}

	return &EncodingTest{
		Name: "encoding",
		Data: byts,
	}
}
