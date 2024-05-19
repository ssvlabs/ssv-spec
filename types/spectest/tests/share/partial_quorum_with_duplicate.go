package share

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialQuorumWithDuplicate tests msg with unique f+1 signers (but also including duplicates)
func PartialQuorumWithDuplicate() *ShareTest {
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
		Signers: []types.OperatorID{1, 1, 2},
	}

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
