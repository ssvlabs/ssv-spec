package postconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func PostConsensusNotDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	partialSigMsg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)
	signedMsg := testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, partialSigMsg))

	return msgvalidation.NewMsgValidationSpecTest(
		"post consensus not decided",
		testdoc.PostConsensusNotDecidedDoc,
		msgvalidation.NewRunnerPreset(
			msgvalidation.RunnerTypeCommittee,
			msgvalidation.StateWithRunningInstanceOp(
				msgvalidation.AttesterDuty(spec.DataVersionPhase0),
				testingutils.TestBeaconVoteByts,
				false,
			),
		),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
