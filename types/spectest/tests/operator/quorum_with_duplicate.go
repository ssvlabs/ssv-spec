package operator

import (
	"crypto/rsa"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// QuorumWithDuplicate tests msg with unique 2f+1 signers (but also including duplicates)
func QuorumWithDuplicate() *OperatorTest {
	ks := testingutils.Testing4SharesSet()
	operator := testingutils.TestingOperator(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[4], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 4, 2, 3})
	msg.OperatorIDs = []types.OperatorID{1, 1, 2, 3}

	return &OperatorTest{
		Name:                  "quorum with duplicate",
		Operator:              *operator,
		Message:               *msg,
		ExpectedHasQuorum:     true,
		ExpectedFullCommittee: false,
		ExpectedError:         "non unique signer",
	}
}
