package general

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func UnsupportedSSVMessageType() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return msgvalidation.NewMsgValidationSpecTest(
		"unsupported ssv message type",
		testdoc.UnsupportedSSVMessageTypeDoc,
		msgvalidation.NewRunnerPreset(msgvalidation.RunnerTypeAggregatorCommittee),
		msgvalidation.EncodeSignedSSVMessage(msgvalidation.MakeUnsupportedSSVMessage(ks)),
		pubsub.ValidationReject,
		ks,
	)
}
