package duty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SyncCommitteeZeroCommitteeMetadata tests that sync committee roles accept zero CommitteeLength and ValidatorCommitteeIndex
func SyncCommitteeZeroCommitteeMetadata() *ValidationTest {
	version := spec.DataVersionElectra
	syncCommittee := testingutils.TestingSyncCommitteeDuty(version).ValidatorDuties[0]
	contribution := testingutils.TestingAggregatorCommitteeDuty(nil, []int{testingutils.TestingValidatorIndex}, version).ValidatorDuties[0]

	duties := []*types.ValidatorDuty{syncCommittee, contribution}
	for _, duty := range duties {
		duty.CommitteeLength = 0
		duty.ValidatorCommitteeIndex = 0
	}

	return NewValidationTest(
		"sync committee zero committee metadata",
		testdoc.ValidatorDutyValidationSyncCommitteeZeroCommitteeMetadataDoc,
		duties,
		0,
	)
}
