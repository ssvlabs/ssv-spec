package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ElectraBlockValidation tests a valid consensus data with electra block
func ElectraBlockValidation() *ProposerConsensusDataTest {
	return NewProposerConsensusDataTest(
		"valid electra block",
		testdoc.ProposerConsensusDataTestElectraBlockDoc,
		*testingutils.TestProposerConsensusDataV(spec.DataVersionElectra),
		0,
	)
}
