package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
)

// PubSubNoData tests an empty data for the pubsub message
func PubSubNoData() tests.SpecTest {
	return &MessageValidationTest{
		Name:          "pubsub no data",
		Messages:      [][]byte{{}},
		ExpectedError: validation.ErrPubSubMessageHasNoData.Error(),
	}
}
