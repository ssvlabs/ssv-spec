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

	// A runner already executing the duty, with its observation frozen. The state is set directly, not by
	// starting the duty here: that would broadcast while the vector is built, which a run from the JSON cannot
	// reproduce.
	running := testingutils.PTCAttesterRunner(ks)
	running.GetBaseRunner().State = ssv.NewRunnerState(ks.Threshold, testingutils.TestingPTCAttesterDuty())
	running.(*ssv.PTCAttesterRunner).PayloadAttestationData = testingutils.TestingPayloadAttestationData(testingutils.TestingDutySlotGloas)

	return newduty.NewMultiStartNewRunnerDutySpecTest(
		"ptc attestation observation",
		testdoc.PTCAttesterObservationDoc,
		[]*newduty.StartNewRunnerDutySpecTest{
			{
				// The post-state pins the observation; the rejected duty signs and broadcasts nothing.
				Name:              "rejected duplicate duty keeps observation",
				Runner:            running,
				Duty:              testingutils.TestingPTCAttesterDuty(),
				Threshold:         ks.Threshold,
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
