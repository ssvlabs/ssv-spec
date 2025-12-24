package aggregatorcommitteeconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SyncCommitteeContributionValidation tests a valid consensus data with sync committee contrib.
func SyncCommitteeContributionValidation() *AggregatorCommitteeConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"sync committee contribution valid",
		testdoc.AggregatorCommitteeConsensusDataTestSyncCommitteeContributionValidationDoc,
		*testingutils.TestSyncCommitteeContributionConsensusData,
		0,
	)
}
