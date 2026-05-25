package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CapellaBlockValidation tests a valid consensus data with capella block
func CapellaBlockValidation() *ProposerConsensusDataTest {
	return NewProposerConsensusDataTest(
		"valid capella block",
		testdoc.ProposerConsensusDataTestCapellaBlockDoc,
		*testingutils.TestProposerConsensusDataV(spec.DataVersionCapella),
		0,
	)
}
