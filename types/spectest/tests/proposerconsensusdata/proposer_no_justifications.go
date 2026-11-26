package proposerconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types/gloas"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ProposerNoJustifications tests an invalid consensus data with no proposer justifications
func ProposerNoJustifications() *ProposerConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	cd := testingutils.TestProposerConsensusDataV(gloas.DataVersionGloas)

	return NewProposerConsensusDataTest(
		"proposer no justification",
		testdoc.ProposerConsensusDataTestProposerNoJustificationsDoc,
		*cd,
		0,
	)
}
