package share

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoPartialQuorumDuplicate tests msg with < unique f+1 signers (but f+1 signers including duplicates)
func NoPartialQuorumDuplicate() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2]}, []types.OperatorID{1, 2})
	msg.OperatorIDs = []types.OperatorID{1, 1}

	return &ShareTest{
		Name:                     "no partial quorum duplicate",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: false,
		ExpectedHasQuorum:        false,
		ExpectedFullCommittee:    false,
		ExpectedError:            "non unique signer",
	}
}
