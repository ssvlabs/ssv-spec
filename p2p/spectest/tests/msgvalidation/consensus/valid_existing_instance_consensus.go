package consensus

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func ValidExistingInstanceConsensus() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	height := qbft.Height(testingutils.TestingDutySlot)

	return msgvalidation.NewMsgValidationSpecTest(
		"valid existing instance consensus",
		testdoc.ValidExistingInstanceConsensusDoc,
		msgvalidation.NewRunnerPreset(
			msgvalidation.RunnerTypeCommittee,
			msgvalidation.StartInstanceOp(height, testingutils.TestBeaconVoteByts),
		),
		msgvalidation.EncodeSignedSSVMessage(msgvalidation.MakeExistingInstanceConsensusMessage(ks, height)),
		pubsub.ValidationAccept,
		ks,
	)
}
