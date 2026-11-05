package ptcattester

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/duties/newduty"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Observation tests the SIP #94 §3 observation lifecycle at duty start: a rejected duplicate duty keeps the running
// duty's frozen observation, and payload attestation data for another slot is rejected without freezing anything.
func Observation() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// A runner already executing the duty, with its observation frozen.
	running := testingutils.PTCAttesterRunner(ks)
	if err := running.StartNewDuty(testingutils.TestingPTCAttesterDuty(), ks.Threshold); err != nil {
		panic(err.Error())
	}
	if running.(*ssv.PTCAttesterRunner).PayloadAttestationData == nil {
		panic("running duty froze no observation")
	}

	return newduty.NewMultiStartNewRunnerDutySpecTest(
		"ptc attestation observation",
		testdoc.PTCAttesterObservationDoc,
		[]*newduty.StartNewRunnerDutySpecTest{
			{
				Name:      "rejected duplicate duty keeps observation",
				Runner:    running,
				Duty:      testingutils.TestingPTCAttesterDuty(),
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusPTCMsg(ks.Shares[1], 1), // broadcast by the running duty
				},
				ExpectedErrorCode: types.DutyAlreadyPassedErrorCode,
			},
			{
				Name:              "wrong slot payload attestation data",
				Runner:            testingutils.PTCAttesterRunner(ks),
				Duty:              testingutils.TestingPTCAttesterDuty(),
				Threshold:         ks.Threshold,
				BeaconNode:        &tests.BeaconNodeBehaviour{WrongPayloadAttestationSlot: true},
				ExpectedErrorCode: types.PTCAttesterWrongSlotErrorCode,
			},
		},
		ks,
	)
}
