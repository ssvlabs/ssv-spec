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

func MultiSignerPartialSignature() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	duty := testingutils.TestingAggregatorAndSyncCommitteeContributorDuties(spec.DataVersionPhase0)
	keySetMap := map[phase0.ValidatorIndex]*testingutils.TestKeySet{
		testingutils.TestingValidatorIndex: ks,
	}
	partialSigMsg := testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, keySetMap, 1)
	signedMsg := testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, partialSigMsg))
	sig2, err := types.SignSSVMessage(ks.OperatorKeys[2], signedMsg.SSVMessage)
	if err != nil {
		panic(err)
	}
	signedMsg.OperatorIDs = append(signedMsg.OperatorIDs, 2)
	signedMsg.Signatures = append(signedMsg.Signatures, sig2)

	return msgvalidation.NewMsgValidationSpecTest(
		"multi signer partial signature",
		testdoc.MultiSignerPartialSignatureDoc,
		msgvalidation.NewRunnerPreset(
			msgvalidation.RunnerTypeAggregatorCommittee,
			msgvalidation.StartDutyOp(msgvalidation.AggregatorCommitteeDuty(spec.DataVersionPhase0)),
		),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
