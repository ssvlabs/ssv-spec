package ptcattester

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DivergingObservation tests honest convergence under a minority split (SIP #94 §3): two peers
// converge on the operator's observation and a third signs a diverging one (same block, different
// payload status). The diverging partial fails the expected-root check, the two matching partials
// stay below quorum, and nothing reconstructs or submits.
func DivergingObservation() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MsgProcessingSpecTest{
		Name:          "ptc attestation diverging observation",
		Documentation: testdoc.PTCAttesterDivergingObservationDoc,
		Runner:        testingutils.PTCAttesterRunner(ks),
		Duty:          testingutils.TestingPTCAttesterDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPTCAttester(nil, testingutils.PreConsensusPTCMsg(ks.Shares[1], 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPTCAttester(nil, testingutils.PreConsensusPTCMsg(ks.Shares[2], 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPTCAttester(nil, testingutils.PreConsensusPTCWrongRootMsg(ks.Shares[3], 3))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusPTCMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		ExpectedErrorCode: types.WrongSigningRootErrorCode,
	}
}
