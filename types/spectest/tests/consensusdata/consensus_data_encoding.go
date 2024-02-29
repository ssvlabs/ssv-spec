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

// Proposer
func ProposerConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("proposer encoding", testingutils.TestProposerConsensusDataV(spec.DataVersionCapella))
}
func ProposerWithJustificationConsensusDataEncoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	return ConsensusDataEncoding("proposer with justification encoding", testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionCapella))
}

// Blinded proposer
func BlindedProposerConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("blinded proposer encoding", testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella))
}
func BlindedProposerWithJustificationConsensusDataEncoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	return ConsensusDataEncoding("blinded proposer with justification encoding", testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionCapella))
}

// Attestation
func AttestationConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("attestation encoding", testingutils.TestAttesterConsensusData)
}
func AttestationWithJustificationConsensusDataEncoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	return ConsensusDataEncoding("attestation with justification encoding", testingutils.TestAttesterWithJustificationsConsensusData(ks))
}

// Aggregator
func AggregatorConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("aggregation encoding", testingutils.TestAggregatorConsensusData)
}
func AggregatorWithJustificationConsensusDataEncoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	return ConsensusDataEncoding("aggregation with justification encoding", testingutils.TestAggregatorWithJustificationsConsensusData(ks))
}

// Sync committee
func SyncCommitteeConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("sync committee encoding", testingutils.TestSyncCommitteeConsensusData)
}
func SyncCommitteeWithJustificationConsensusDataEncoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	return ConsensusDataEncoding("sync committee with justification encoding", testingutils.TestSyncCommitteeWithJustificationsConsensusData(ks))
}

// Sync committee contribution
func SyncCommitteeContributionConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding("sync committee contribution encoding", testingutils.TestSyncCommitteeContributionConsensusData)
}
func SyncCommitteeWithJustificationContributionConsensusDataEncoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	return ConsensusDataEncoding("sync committee contribution with justification encoding", testingutils.TestSyncCommitteeContributionWithJustificationConsensusData(ks))
}
