package share

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoQuorumDuplicate tests msg with < unique 2f+1 signers (but 2f+1 signers including duplicates)
func NoQuorumDuplicate() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[1], ks.NetworkKeys[2]}, []types.OperatorID{1, 1, 2})

	return &ShareTest{
		Name:                     "no quorum duplicate",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: true,
		ExpectedHasQuorum:        false,
		ExpectedFullCommittee:    false,
		ExpectedError:            "non unique signer",
	}
}
