package consensus

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func ValidFutureDecidedConsensus() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return msgvalidation.NewMsgValidationSpecTest(
		"valid future decided consensus",
		testdoc.ValidFutureDecidedConsensusDoc,
		msgvalidation.NewRunnerPreset(msgvalidation.RunnerTypeCommittee),
		msgvalidation.EncodeSignedSSVMessage(msgvalidation.MakeFutureDecidedConsensusMessage(ks)),
		pubsub.ValidationAccept,
		ks,
	)
}
