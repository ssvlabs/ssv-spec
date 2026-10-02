package committeemember

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidOperator tests committee members with an operator that fails Operator.Validate (zero ID, nil)
func InvalidOperator() *ValidationTest {
	ks := testingutils.Testing4SharesSet()

	zeroID := testingutils.TestingCommitteeMember(ks)
	zeroID.Committee[3].OperatorID = 0
	// Matching CommitteeID, so only the operator check fails
	zeroID.CommitteeID = types.GetCommitteeID([]types.OperatorID{1, 2, 3, 0})

	nilOperator := testingutils.TestingCommitteeMember(ks)
	nilOperator.Committee[3] = nil

	return NewValidationTest(
		"invalid operator",
		testdoc.CommitteeMemberValidationInvalidOperatorDoc,
		[]*types.CommitteeMember{zeroID, nilOperator},
		types.InvalidCommitteeMemberErrorCode,
	)
}
