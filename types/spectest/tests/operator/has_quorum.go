package operator

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HasQuorum tests msg with unique 2f+1 signers
func HasQuorum() *OperatorTest {
	ks := testingutils.Testing4SharesSet()
	operator := testingutils.TestingOperator(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3})

	return &OperatorTest{
		Name:                  "has quorum",
		Operator:              *operator,
		Message:               *msg,
		ExpectedHasQuorum:     true,
		ExpectedFullCommittee: false,
	}
}
