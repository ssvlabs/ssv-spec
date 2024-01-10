package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Valid tests a valid SignedSSVMessageTest
func Valid() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingSignedSSVMessage(ks.Shares[1], 1)

	return &SignedSSVMessageTest{
		Name: "valid",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
	}
}
