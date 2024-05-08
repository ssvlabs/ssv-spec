package attestationsynccommittee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func validatorIndexList(limit int) []int {
	ret := make([]int, limit)
	for i := 0; i < limit; i++ {
		ret[i] = i + 1
	}
	return ret
}

// StartAttestations starts a cluster runner with a list of 30 attestation duties
func StartAttestations() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "start attestations",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:           "30 attestations",
				Runner:         testingutils.CommitteeRunner(ks),
				Duty:           testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot, validatorIndexList(30)),
				Messages:       []*types.SignedSSVMessage{},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
