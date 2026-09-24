package consensus

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func UnknownInstanceConsensus() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	height := qbft.Height(testingutils.TestingDutySlot)

	return msgvalidation.NewMsgValidationSpecTest(
		"unknown instance consensus",
		testdoc.UnknownInstanceConsensusDoc,
		msgvalidation.NewRunnerPreset(
			msgvalidation.RunnerTypeCommittee,
			msgvalidation.SetControllerHeightOp(height),
		),
		msgvalidation.EncodeSignedSSVMessage(msgvalidation.MakeExistingInstanceConsensusMessage(ks, height)),
		pubsub.ValidationReject,
		ks,
	)
}
