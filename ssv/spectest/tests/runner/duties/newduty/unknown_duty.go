package newduty

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownDuty tests an invalid UnknownDuty
func UnknownDuty() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()
	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty unknown duty",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "unknown duty",
				Runner:                  testingutils.UnknownDutyTypeRunner(ks),
				Duty:                    testingutils.TestingUnknownDutyType,
				PostDutyRunnerStateRoot: "74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "runner beacon role unknown",
			},
		},
	}
}
