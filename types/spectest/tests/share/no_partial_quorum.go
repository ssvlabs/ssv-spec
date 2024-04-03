package share

import (
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoPartialQuorum tests  msg with < unique f+1 signers
func NoPartialQuorum() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1)

	return &ShareTest{
		Name:                     "no partial quorum",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: false,
		ExpectedHasQuorum:        false,
		ExpectedFullCommittee:    false,
	}
}
