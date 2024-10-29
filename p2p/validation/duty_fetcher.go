package validation

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Check duty assignments
type DutyFetcher interface {
	HasProposerDuty(validatorIndex phase0.ValidatorIndex, slot phase0.Slot) bool
	HasSyncCommitteeContributionDuty(validatorIndex phase0.ValidatorIndex, slot phase0.Slot) bool
	HasSyncCommitteeDuty(validatorIndex phase0.ValidatorIndex, slot phase0.Slot) bool
}
