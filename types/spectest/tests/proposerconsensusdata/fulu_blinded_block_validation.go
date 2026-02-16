package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FuluBlindedBlockValidation tests a valid consensus data with fulu blinded block
func FuluBlindedBlockValidation() *ProposerConsensusDataTest {
	return NewProposerConsensusDataTest(
		"valid fulu blinded block",
		testdoc.ProposerConsensusDataTestElectraBlindedBlockDoc,
		*testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionFulu),
		0,
	)
}
