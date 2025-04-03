package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ProposerNoJustifications tests an invalid consensus data with no proposer justifications
func ProposerNoJustifications() *ValidatorConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	cd := testingutils.TestProposerConsensusDataV(spec.DataVersionCapella)

	return &ValidatorConsensusDataTest{
		Name:          "proposer no justification",
		ConsensusData: *cd,
	}
}
