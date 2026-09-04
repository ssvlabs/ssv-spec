package committeesingleduty

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasPreGloasVote sends a proposal carrying a pre-Gloas BeaconVote to a committee running a duty at
// a Gloas slot. The Gloas consensus value is a fixed 120-byte GloasBeaconVote (SIP #94 §2), so the
// 112-byte pre-Gloas value must be rejected by the value check.
//
// This is the runner-level counterpart to the gloas attestation value-check vectors, and it pins the
// rule that the fork is selected by the duty's slot rather than by the shape of the proposed value:
// dispatching on shape would accept this value here and only fail later, when the runner decodes the
// decided value — killing the duty after it had already decided.
func GloasPreGloasVote() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	ksMap := testingutils.KeySetMapForValidators(1)
	gloasSlot := phase0.Slot(testingutils.TestingDutySlotGloas)
	gloasVersion := testingutils.VersionBySlot(gloasSlot)

	return committee.NewMultiCommitteeSpecTest(
		"gloas pre-gloas vote",
		testdoc.CommitteeGloasPreGloasVoteDoc,
		[]*committee.CommitteeSpecTest{
			{
				Name:      "1 attestation",
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingAttesterDutyForValidators(gloasVersion, testingutils.ValidatorIndexList(1)),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), testingutils.CommitteeMsgID(ks),
						testingutils.TestBeaconVoteByts, // pre-Gloas 112-byte value at a Gloas slot
						qbft.Height(gloasSlot)),
				},
				ExpectedErrorCode: types.DecodeGloasBeaconVoteErrorCode,
			},
		},
		ks,
	)
}
