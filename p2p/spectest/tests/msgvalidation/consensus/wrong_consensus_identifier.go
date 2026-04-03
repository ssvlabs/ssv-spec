package consensus

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func WrongConsensusIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	height := qbft.Height(testingutils.TestingDutySlot)
	signedMsg := testingutils.TestingProposalMessageWithIdentifierAndFullData(
		ks.OperatorKeys[1],
		1,
		testingutils.AggregatorCommitteeMsgID(ks),
		testingutils.TestBeaconVoteByts,
		height,
	)

	return msgvalidation.NewMsgValidationSpecTest(
		"wrong consensus identifier",
		testdoc.WrongConsensusIdentifierDoc,
		msgvalidation.NewRunnerPreset(
			msgvalidation.RunnerTypeCommittee,
			msgvalidation.StartInstanceOp(height, testingutils.TestBeaconVoteByts),
		),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
