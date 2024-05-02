package share

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// QuorumWithDuplicate tests msg with unique 2f+1 signers (but also including duplicates)
func QuorumWithDuplicate() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := &qbft.SignedMessage{
		Message: qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.TestingIdentifier,
			Root:       testingutils.TestingQBFTRootData,
		},
		Signers: []types.OperatorID{1, 1, 2, 3},
	}

	return &ShareTest{
		Name:                     "quorum with duplicate",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: true,
		ExpectedHasQuorum:        true,
		ExpectedFullCommittee:    false,
		ExpectedError:            "non unique signer",
	}
}
