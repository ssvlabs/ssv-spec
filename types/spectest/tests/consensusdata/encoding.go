package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Encoding tests encoding of a ConsensusData struct
func Encoding() *EncodingSpecTest {
	data := &types.ConsensusData{
		Duty:                   testingutils.TestingAttesterDuty,
		AttestationData:        testingutils.TestingAttestationData,
		BlockData:              testingutils.TestingBeaconBlock,
		AggregateAndProof:      testingutils.TestingAggregateAndProof,
		SyncCommitteeBlockRoot: testingutils.TestingSyncCommitteeBlockRoot,
		SyncCommitteeContribution: types.ContributionsMap{
			testingutils.TestingContributionProofsSigned[0]: testingutils.TestingSyncCommitteeContributions[0],
			testingutils.TestingContributionProofsSigned[1]: testingutils.TestingSyncCommitteeContributions[1],
			testingutils.TestingContributionProofsSigned[2]: testingutils.TestingSyncCommitteeContributions[2],
		},
	}

	return &EncodingSpecTest{
		Name: "encoding ConsensusData",
		Obj:  data,
	}
}
