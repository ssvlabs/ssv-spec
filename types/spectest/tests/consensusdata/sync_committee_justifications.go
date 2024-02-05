package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SyncCommitteeJustifications tests an invalid consensus data with sync committee pre-consensus justifications
func SyncCommitteeJustifications() *ConsensusDataTest {
	ks := testingutils.Testing4SharesSet()

	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(0); i <= ks.Threshold; i++ {
		justif = append(justif, testingutils.PreConsensusRandaoMsg(ks.Shares[i+1], i+1))
	}

	cd := types.ConsensusData{
		Duty:                       testingutils.TestingSyncCommitteeDuty,
		PreConsensusJustifications: justif,
		Version:                    spec.DataVersionCapella,
		DataSSZ:                    testingutils.TestingSyncCommitteeBlockRoot[:],
	}

	return &ConsensusDataTest{
		Name:          "sync committee with pre-consensus justification",
		ConsensusData: cd,
		ExpectedError: "sync committee invalid justifications",
	}
}
