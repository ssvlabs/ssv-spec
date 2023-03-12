package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Encoding tests encoding of a ssv message
func Encoding() *EncodingTest {
	msg := testingutils.TestSyncCommitteeContributionConsensusData

	byts, err := msg.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := msg.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &EncodingTest{
		Name:         "ConsensusData encoding",
		Data:         byts,
		ExpectedRoot: root,
	}
}
