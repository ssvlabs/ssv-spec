package operator

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HasQuorum3f1 tests msg with unique 3f+1 signers
func HasQuorum3f1() *OperatorTest {
	ks := testingutils.Testing4SharesSet()
	operator := testingutils.TestingOperator(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3], ks.OperatorKeys[4]}, []types.OperatorID{1, 2, 3, 4})

	return &OperatorTest{
		Name:                  "has quorum 3f1",
		Operator:              *operator,
		Message:               *msg,
		ExpectedHasQuorum:     true,
		ExpectedFullCommittee: true,
	}
}
