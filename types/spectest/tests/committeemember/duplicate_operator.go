package committeemember

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicateOperator tests a committee member whose committee repeats an operator ID
func DuplicateOperator() *ValidationTest {
	cm := testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet())
	cm.Committee[1].OperatorID = cm.Committee[0].OperatorID
	// Matching CommitteeID, so only the duplicate check fails
	cm.CommitteeID = types.GetCommitteeID([]types.OperatorID{1, 1, 3, 4})

	return NewValidationTest(
		"duplicate operator",
		testdoc.CommitteeMemberValidationDuplicateOperatorDoc,
		[]*types.CommitteeMember{cm},
		types.InvalidCommitteeMemberErrorCode,
	)
}
