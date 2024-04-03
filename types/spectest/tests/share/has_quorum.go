package share

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HasQuorum tests msg with unique 2f+1 signers
func HasQuorum() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]}, []types.OperatorID{1, 2, 3})

	return &ShareTest{
		Name:                     "has quorum",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: true,
		ExpectedHasQuorum:        true,
		ExpectedFullCommittee:    false,
	}
}
