package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ProposerNoJustifications tests an invalid consensus data with no proposer justifications
func ProposerNoJustifications() *ProposerConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	cd := testingutils.TestProposerConsensusDataV(spec.DataVersionCapella)

	return NewProposerConsensusDataTest(
		"proposer no justification",
		testdoc.ProposerConsensusDataTestProposerNoJustificationsDoc,
		*cd,
		0,
	)
}
