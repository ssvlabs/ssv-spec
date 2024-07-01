package committeemember

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HasQuorum tests msg with unique 2f+1 signers
func HasQuorum() *CommitteeMemberTest {
	ks := testingutils.Testing4SharesSet()
	committeeMember := testingutils.TestingCommitteeMember(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3})

	return &CommitteeMemberTest{
		Name:                  "has quorum",
		CommitteeMember:       *committeeMember,
		Message:               *msg,
		ExpectedHasQuorum:     true,
		ExpectedFullCommittee: false,
	}
}
