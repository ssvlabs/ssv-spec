package attestationsynccommittee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartSyncCommittees starts a cluster runner with a list of 30 sync committee duties
func StartSyncCommittees() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "start sync committees",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:           "30 sync committees",
				Runner:         testingutils.CommitteeRunner(ks),
				Duty:           testingutils.TestingCommitteeSyncCommitteeDuty(testingutils.TestingDutySlot, validatorIndexList(30)),
				Messages:       []*types.SignedSSVMessage{},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
