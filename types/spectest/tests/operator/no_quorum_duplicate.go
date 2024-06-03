package operator

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoQuorumDuplicate tests msg with < unique 2f+1 signers (but 2f+1 signers including duplicates)
func NoQuorumDuplicate() *OperatorTest {
	ks := testingutils.Testing4SharesSet()
	operator := testingutils.TestingOperator(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[3], ks.OperatorKeys[2]}, []types.OperatorID{1, 3, 2})
	msg.OperatorIDs = []types.OperatorID{1, 1, 2}

	return &OperatorTest{
		Name:                  "no quorum duplicate",
		Operator:              *operator,
		Message:               *msg,
		ExpectedHasQuorum:     false,
		ExpectedFullCommittee: false,
		ExpectedError:         "non unique signer",
	}
}
