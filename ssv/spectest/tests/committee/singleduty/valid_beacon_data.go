package committeesingleduty

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidBeaconVote sends a proposal message to a maximal committee runner with a valid BeaconVote
func ValidBeaconVote() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	ksMapFor30Validators := testingutils.KeySetMapForValidators(30)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "valid beacon vote",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &committee.CommitteeSpecTest{
			Name:      fmt.Sprintf("30 attestations 30 sync committees (%s)", version.String()),
			Committee: testingutils.BaseCommittee(ksMapFor30Validators),
			Input: []interface{}{
				testingutils.TestingCommitteeDuty(testingutils.ValidatorIndexList(30), testingutils.ValidatorIndexList(30), version),
				testingutils.TestingProposalMessageWithIdentifierAndFullData(
					ks.OperatorKeys[1], types.OperatorID(1), testingutils.CommitteeMsgID(ks), testingutils.TestBeaconVoteByts,
					qbft.Height(testingutils.TestingDutySlotV(version))),
			},
			OutputMessages: []*types.PartialSignatureMessages{},
		})
	}

	return multiSpecTest
}
