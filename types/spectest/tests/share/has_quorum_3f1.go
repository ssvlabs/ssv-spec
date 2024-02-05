package share

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// HasQuorum3f1 tests msg with unique 3f+1 signers
func HasQuorum3f1() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3], ks.Shares[4]}, []types.OperatorID{1, 2, 3, 4})

	return &ShareTest{
		Name:                     "has quorum 3f1",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: true,
		ExpectedHasQuorum:        true,
		ExpectedFullCommittee:    true,
	}
}
