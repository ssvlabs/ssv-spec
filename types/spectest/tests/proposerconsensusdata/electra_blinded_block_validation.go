package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ElectraBlindedBlockValidation tests a valid consensus data with electra blinded block
func ElectraBlindedBlockValidation() *ProposerConsensusDataTest {
	return NewProposerConsensusDataTest(
		"valid electra blinded block",
		testdoc.ProposerConsensusDataTestElectraBlindedBlockDoc,
		*testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionElectra),
		0,
	)
}
