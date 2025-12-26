package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FuluBlockValidation tests a valid consensus data with fulu block
func FuluBlockValidation() *ProposerConsensusDataTest {
	return NewProposerConsensusDataTest(
		"valid fulu block",
		testdoc.ProposerConsensusDataTestElectraBlockDoc,
		*testingutils.TestProposerConsensusDataV(spec.DataVersionFulu),
		0,
	)
}
