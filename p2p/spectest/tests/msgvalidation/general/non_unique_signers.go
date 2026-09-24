package general

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

func NonUniqueSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	duty := testingutils.TestingAggregatorAndSyncCommitteeContributorDuties(spec.DataVersionPhase0)
	keySetMap := map[phase0.ValidatorIndex]*testingutils.TestKeySet{
		testingutils.TestingValidatorIndex: ks,
	}
	partialSigMsg := testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, keySetMap, 1)
	base := testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, partialSigMsg))
	signedMsg := msgvalidation.MutateSignedSSVMessage(base, func(msg *types.SignedSSVMessage) {
		msg.OperatorIDs = append(msg.OperatorIDs, msg.OperatorIDs[0])
		msg.Signatures = append(msg.Signatures, msg.Signatures[0])
	})

	return msgvalidation.NewMsgValidationSpecTest(
		"non unique signers",
		testdoc.NonUniqueSignersDoc,
		msgvalidation.NewRunnerPreset(msgvalidation.RunnerTypeAggregatorCommittee),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
