package share

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HasQuorum3f1 tests msg with unique 3f+1 signers
func HasQuorum3f1() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3], ks.NetworkKeys[4]}, []types.OperatorID{1, 2, 3, 4})

	return &ShareTest{
		Name:                     "has quorum 3f1",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: true,
		ExpectedHasQuorum:        true,
		ExpectedFullCommittee:    true,
	}
}
