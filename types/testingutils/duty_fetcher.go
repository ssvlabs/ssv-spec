package testingutils

import "github.com/attestantio/go-eth2-client/spec/phase0"

// Mock values
var (
	ValidatorIndexWithoutProposerDuty                  = phase0.ValidatorIndex(1234)
	ValidatorIndexWithoutSyncCommitteeContributionDuty = phase0.ValidatorIndex(5678)
	ValidatorIndexWithSyncCommitteeDuty                = phase0.ValidatorIndex(91011)
)

type TestingDutyFetcher struct {
}

func NewTestingDutyFetcher() *TestingDutyFetcher {
	return &TestingDutyFetcher{}
}

func (df *TestingDutyFetcher) HasProposerDuty(validatorIndex phase0.ValidatorIndex, slot phase0.Slot) bool {
	return (ValidatorIndexWithoutProposerDuty != validatorIndex)
}
func (df *TestingDutyFetcher) HasSyncCommitteeContributionDuty(validatorIndex phase0.ValidatorIndex, slot phase0.Slot) bool {
	return (ValidatorIndexWithoutSyncCommitteeContributionDuty != validatorIndex)
}
func (df *TestingDutyFetcher) HasSyncCommitteeDuty(validatorIndex phase0.ValidatorIndex, slot phase0.Slot) bool {
	return (ValidatorIndexWithSyncCommitteeDuty == validatorIndex)
}
