package consensus

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func InvalidConsensusSignature() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	height := qbft.Height(testingutils.TestingDutySlot)
	signedMsg := msgvalidation.MakeExistingInstanceConsensusMessage(ks, height)
	signedMsg.OperatorIDs[0] = 2

	return msgvalidation.NewMsgValidationSpecTest(
		"invalid consensus signature",
		testdoc.InvalidConsensusSignatureDoc,
		msgvalidation.NewRunnerPreset(
			msgvalidation.RunnerTypeCommittee,
			msgvalidation.StartInstanceOp(height, testingutils.TestBeaconVoteByts),
		),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
