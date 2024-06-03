package committeemember

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoQuorumDuplicate tests msg with < unique 2f+1 signers (but 2f+1 signers including duplicates)
func NoQuorumDuplicate() *CommitteeMemberTest {
	ks := testingutils.Testing4SharesSet()
	committeeMember := testingutils.TestingCommitteeMember(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[3], ks.OperatorKeys[2]}, []types.OperatorID{1, 3, 2})
	msg.OperatorIDs = []types.OperatorID{1, 1, 2}

	return &CommitteeMemberTest{
		Name:                  "no quorum duplicate",
		CommitteeMember:       *committeeMember,
		Message:               *msg,
		ExpectedHasQuorum:     false,
		ExpectedFullCommittee: false,
		ExpectedError:         "non unique signer",
	}
}
