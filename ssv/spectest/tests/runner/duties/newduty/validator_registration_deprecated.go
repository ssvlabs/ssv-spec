package newduty

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidatorRegistrationDeprecated tests that the validator registration duty is rejected from Gloas
// (SIP #94 §5): fee recipient and gas limit travel in the proposer preferences instead. Pre-Gloas
// registrations are unchanged (covered by every other validator registration vector).
func ValidatorRegistrationDeprecated() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	gloasDuty := testingutils.TestingValidatorRegistrationDuty
	gloasDuty.Slot = testingutils.TestingDutySlotGloas

	return NewMultiStartNewRunnerDutySpecTest(
		"new duty validator registration deprecated",
		testdoc.NewDutyValidatorRegistrationDeprecatedDoc,
		[]*StartNewRunnerDutySpecTest{
			{
				Name:                    "validator registration at a gloas slot",
				Runner:                  testingutils.ValidatorRegistrationRunner(ks),
				Duty:                    &gloasDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				ExpectedErrorCode:       types.ValidatorRegistrationDeprecatedErrorCode,
			},
		},
		ks,
	)
}
