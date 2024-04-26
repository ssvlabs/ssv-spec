package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ConsensusDataEncoding tests encoding and decoding ConsensusData for all duties
func ConsensusDataEncoding(name string, cd *types.ConsensusData) *EncodingTest {

	byts, err := cd.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := cd.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &EncodingTest{
		Name:         name,
		Data:         byts,
		ExpectedRoot: root,
	}
}

func ProposerConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("proposer encoding", testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella))
}
func BlindedProposerConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("blinded proposer encoding", testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella))
}
func AggregatorConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("aggregation encoding", testingutils.TestAggregatorConsensusData)
}
func SyncCommitteeContributionConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("sync committee contribution encoding", testingutils.TestSyncCommitteeContributionConsensusData)
}
