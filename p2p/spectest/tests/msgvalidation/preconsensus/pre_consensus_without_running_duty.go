package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func PreConsensusWithoutRunningDuty() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	duty := testingutils.TestingAggregatorAndSyncCommitteeContributorDuties(spec.DataVersionPhase0)
	keySetMap := map[phase0.ValidatorIndex]*testingutils.TestKeySet{
		testingutils.TestingValidatorIndex: ks,
	}
	partialSigMsg := testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, keySetMap, 1)
	signedMsg := testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, partialSigMsg))

	return msgvalidation.NewMsgValidationSpecTest(
		"pre consensus without running duty",
		testdoc.PreConsensusWithoutRunningDutyDoc,
		msgvalidation.NewRunnerPreset(msgvalidation.RunnerTypeAggregatorCommittee),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
