package share

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialQuorumWithDuplicate tests msg with unique f+1 signers (but also including duplicates)
func PartialQuorumWithDuplicate() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[3], ks.OperatorKeys[2]}, []types.OperatorID{1, 3, 2})
	msg.OperatorIDs = []types.OperatorID{1, 1, 2}

	return &ShareTest{
		Name:                     "partial quorum with duplicate",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: true,
		ExpectedHasQuorum:        false,
		ExpectedFullCommittee:    false,
		ExpectedError:            "non unique signer",
	}
}
