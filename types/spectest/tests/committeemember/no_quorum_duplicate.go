package committeemember

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoQuorumDuplicate tests msg with < unique 2f+1 signers (but 2f+1 signers including duplicates)
func NoQuorumDuplicate() *CommitteeMemberTest {
	ks := testingutils.Testing4SharesSet()
	committeeMember := testingutils.TestingCommitteeMember(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[3], ks.OperatorKeys[2]}, []types.OperatorID{1, 3, 2})
	msg.OperatorIDs = []types.OperatorID{1, 1, 2}

	return NewCommitteeMemberTest(
		"no quorum duplicate",
		testdoc.NoQuorumDuplicateTestDoc,
		*committeeMember,
		*msg,
		false,
		false,
		types.NonUniqueSignerErrorCode,
	)
}
