package committeemember

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HasQuorum3f1 tests msg with unique 3f+1 signers
func HasQuorum3f1() *CommitteeMemberTest {
	ks := testingutils.Testing4SharesSet()
	committeeMember := testingutils.TestingCommitteeMember(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3], ks.OperatorKeys[4]}, []types.OperatorID{1, 2, 3, 4})

	return &CommitteeMemberTest{
		Name:                  "has quorum 3f1",
		CommitteeMember:       *committeeMember,
		Message:               *msg,
		ExpectedHasQuorum:     true,
		ExpectedFullCommittee: true,
	}
}
