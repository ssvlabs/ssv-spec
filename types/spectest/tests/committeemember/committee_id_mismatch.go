package committeemember

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CommitteeIDMismatch tests a committee member whose CommitteeID doesn't match its committee's operator IDs
func CommitteeIDMismatch() *ValidationTest {
	cm := testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet())
	cm.CommitteeID = types.CommitteeID{}

	return NewValidationTest(
		"committee id mismatch",
		testdoc.CommitteeMemberValidationCommitteeIDMismatchDoc,
		[]*types.CommitteeMember{cm},
		types.InvalidCommitteeMemberErrorCode,
	)
}
