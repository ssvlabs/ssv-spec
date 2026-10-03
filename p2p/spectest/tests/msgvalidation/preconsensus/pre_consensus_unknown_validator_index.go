package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func PreConsensusUnknownValidatorIndex() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	base := testingutils.SignPartialSigSSVMessage(
		ks,
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)),
	)
	signedMsg := msgvalidation.MutatePartialSignatureSignedMessage(base, 1, ks.OperatorKeys[1], func(msg *types.PartialSignatureMessages) {
		for _, inner := range msg.Messages {
			inner.ValidatorIndex = testingutils.TestingWrongValidatorIndex
		}
	})

	return msgvalidation.NewMsgValidationSpecTest(
		"pre consensus unknown validator index",
		testdoc.PreConsensusUnknownValidatorIndexDoc,
		msgvalidation.NewRunnerPreset(
			msgvalidation.RunnerTypeProposer,
			msgvalidation.StartDutyOp(msgvalidation.ProposerDuty(spec.DataVersionDeneb)),
		),
		msgvalidation.EncodeSignedSSVMessage(signedMsg),
		pubsub.ValidationReject,
		ks,
	)
}
