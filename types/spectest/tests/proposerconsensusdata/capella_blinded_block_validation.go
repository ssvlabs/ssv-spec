package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CapellaBlindedBlockValidation tests a valid consensus data with capella blinded block
func CapellaBlindedBlockValidation() *ProposerConsensusDataTest {
	return NewProposerConsensusDataTest(
		"valid capella blinded block",
		testdoc.ProposerConsensusDataTestCapellaBlindedBlockDoc,
		*testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella),
		0,
	)
}
