package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ConsensusDataEncoding tests encoding and decoding ValidatorConsensusData for all duties
func ConsensusDataEncoding(name string, cd *types.ValidatorConsensusData) *EncodingTest {

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
func Phase0AggregatorConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("phase0 aggregation encoding", testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0))
}
func ElectraAggregatorConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("electra aggregation encoding", testingutils.TestAggregatorConsensusData(spec.DataVersionElectra))
}
func SyncCommitteeContributionConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("sync committee contribution encoding", testingutils.TestSyncCommitteeContributionConsensusData)
}
