package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func PreConsensusFutureSlot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	duty := testingutils.TestingAggregatorAndSyncCommitteeContributorDuties(spec.DataVersionPhase0)
	keySetMap := map[phase0.ValidatorIndex]*testingutils.TestKeySet{
		testingutils.TestingValidatorIndex: ks,
	}
	base := testingutils.SignPartialSigSSVMessage(
		ks,
		testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, keySetMap, 1)),
	)
	signedMsg := msgvalidation.MutatePartialSignatureSignedMessage(base, 1, ks.OperatorKeys[1], func(msg *types.PartialSignatureMessages) {
		msg.Slot++
	})

	return msgvalidation.NewMsgValidationSpecTest(
		"pre consensus future slot",
		testdoc.PreConsensusFutureSlotDoc,
		msgvalidation.NewRunnerPreset(
			msgvalidation.RunnerTypeAggregatorCommittee,
			msgvalidation.StartDutyOp(msgvalidation.AggregatorCommitteeDuty(spec.DataVersionPhase0)),
		),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
