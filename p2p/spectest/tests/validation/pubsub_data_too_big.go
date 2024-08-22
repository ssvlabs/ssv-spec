package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
)

// PubSubDataTooBig tests a pubsub message with a data size above the limit
func PubSubDataTooBig() tests.SpecTest {

	bigMsg := [validation.MaxEncodedMsgSize + 1]byte{}

	return &MessageValidationTest{
		Name:          "pubsub data too big",
		Messages:      [][]byte{bigMsg[:]},
		ExpectedError: validation.ErrPubSubDataTooBig.Error(),
	}
}
