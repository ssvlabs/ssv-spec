package general

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func MalformedPubsubPayload() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return msgvalidation.NewMsgValidationSpecTest(
		"malformed pubsub payload",
		testdoc.MalformedPubsubPayloadDoc,
		msgvalidation.NewRunnerPreset(msgvalidation.RunnerTypeCommittee),
		[]byte{1, 2, 3},
		pubsub.ValidationReject,
		ks,
	)
}
