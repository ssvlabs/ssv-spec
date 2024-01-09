package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// AttesterJustifications tests an invalid consensus data with attester pre-consensus justifications
func AttesterJustifications() *ConsensusDataTest {

	ks := testingutils.Testing4SharesSet()

	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(0); i <= ks.Threshold; i++ {
		justif = append(justif, testingutils.PreConsensusRandaoMsg(ks.Shares[i+1], i+1))
	}

	cd := testingutils.TestAttesterConsensusData

	cd.PreConsensusJustifications = justif

	return &ConsensusDataTest{
		Name:          "invalid attestation with justification",
		ConsensusData: *cd,
		ExpectedError: "attester invalid justifications",
	}
}
