package share

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HasQuorum3f1 tests msg with unique 3f+1 signers
func HasQuorum3f1() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3], ks.OperatorKeys[4]}, []types.OperatorID{1, 2, 3, 4})

	return &ShareTest{
		Name:                     "has quorum 3f1",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: true,
		ExpectedHasQuorum:        true,
		ExpectedFullCommittee:    true,
	}
}
