package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DenebBlockValidation tests a valid consensus data with deneb block
func DenebBlockValidation() *ProposerConsensusDataTest {
	return NewProposerConsensusDataTest(
		"valid deneb block",
		testdoc.ProposerConsensusDataTestDenebBlockDoc,
		*testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb),
		0,
	)
}
