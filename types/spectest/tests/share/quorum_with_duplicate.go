package share

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// QuorumWithDuplicate tests msg with unique 2f+1 signers (but also including duplicates)
func QuorumWithDuplicate() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[4], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 4, 2, 3})
	msg.OperatorIDs = []types.OperatorID{1, 1, 2, 3}

	return &ShareTest{
		Name:                  "quorum with duplicate",
		Share:                 *share,
		Message:               *msg,
		ExpectedHasQuorum:     true,
		ExpectedFullCommittee: false,
		ExpectedError:         "non unique signer",
	}
}
