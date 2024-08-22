package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
)

// PubSubMalformed tests a pubsub message with undecodable data
func PubSubMalformed() tests.SpecTest {

	return &MessageValidationTest{
		Name:          "pubsub malformed",
		Messages:      [][]byte{[]byte{1, 2, 3, 4}},
		ExpectedError: validation.ErrMalformedPubSubMessage.Error(),
	}
}
