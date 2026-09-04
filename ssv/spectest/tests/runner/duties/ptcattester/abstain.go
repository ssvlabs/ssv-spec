package ptcattester

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Abstain tests the SIP #94 §3 abstain rule: the beacon node reports no block for the slot (a zero
// BeaconBlockRoot), so the operator freezes, signs and broadcasts nothing, and with no frozen
// observation every incoming peer partial signature is rejected. The duty stays running: a
// re-triggered duty would re-observe from scratch.
func Abstain() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MsgProcessingSpecTest{
		Name:          "ptc attestation abstain",
		Documentation: testdoc.PTCAttesterAbstainDoc,
		Runner:        testingutils.PTCAttesterRunner(ks),
		Duty:          testingutils.TestingPTCAttesterAbstainDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPTCAttester(nil, testingutils.PreConsensusPTCAbstainSlotMsg(ks.Shares[1], 1))),
		},
		OutputMessages:    []*types.PartialSignatureMessages{},
		ExpectedErrorCode: types.PTCAttesterNoObservationErrorCode,
	}
}
