package committeemember

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FaultyNodesMismatch tests a committee member whose FaultyNodes doesn't satisfy n == 3f+1
func FaultyNodesMismatch() *ValidationTest {
	cm := testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet())
	cm.FaultyNodes = 0

	return NewValidationTest(
		"faulty nodes mismatch",
		testdoc.CommitteeMemberValidationFaultyNodesMismatchDoc,
		[]*types.CommitteeMember{cm},
		types.InvalidCommitteeMemberErrorCode,
	)
}
