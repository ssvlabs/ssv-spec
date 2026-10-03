package consensus

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func DecidedConsensusBadFullData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	signedMsg := msgvalidation.MakeFutureDecidedConsensusMessage(ks)
	signedMsg.FullData = testingutils.TestWrongBeaconVoteByts

	return msgvalidation.NewMsgValidationSpecTest(
		"decided consensus bad full data",
		testdoc.DecidedConsensusBadFullDataDoc,
		msgvalidation.NewRunnerPreset(msgvalidation.RunnerTypeCommittee),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
