package committeesingleduty

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongBeaconVote sends a proposal message to a cluster runner with an invalid beacon vote
func WrongBeaconVote() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	ksMapFor30Validators := testingutils.KeySetMapForValidators(30)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "wrong beacon vote",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &committee.CommitteeSpecTest{
			Name:      fmt.Sprintf("30 attestations 30 sync committees (%s)", version.String()),
			Committee: testingutils.BaseCommittee(ksMapFor30Validators),
			Input: []interface{}{
				testingutils.TestingCommitteeDuty(testingutils.ValidatorIndexList(30), testingutils.ValidatorIndexList(30), version),
				testingutils.TestingProposalMessageWithIdentifierAndFullData(
					ks.OperatorKeys[1], types.OperatorID(1), testingutils.CommitteeMsgID(ks), testingutils.TestWrongBeaconVoteByts,
					qbft.Height(testingutils.TestingDutySlotV(version))),
			},
			OutputMessages: []*types.PartialSignatureMessages{},
			ExpectedError:  "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: attestation data source >= target",
		})
	}

	return multiSpecTest
}
