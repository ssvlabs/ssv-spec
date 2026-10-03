package consensus

import (
	"crypto/rsa"

	"github.com/attestantio/go-eth2-client/spec"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func FutureConsensusMultiSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	signedMsg := testingutils.TestingMultiSignerProposalMessageWithParams(
		[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2]},
		[]types.OperatorID{1, 2},
		qbft.FirstRound,
		qbft.Height(testingutils.TestingDutySlot),
		testingutils.CommitteeMsgID(ks),
		testingutils.TestBeaconVoteByts,
		testingutils.TestingAttestationDataRoot(spec.DataVersionPhase0),
	)

	return msgvalidation.NewMsgValidationSpecTest(
		"future consensus multi signer",
		testdoc.FutureConsensusMultiSignerDoc,
		msgvalidation.NewRunnerPreset(msgvalidation.RunnerTypeCommittee),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
