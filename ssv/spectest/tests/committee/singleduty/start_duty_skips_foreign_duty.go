package committeesingleduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartDutySkipsForeignDuty starts a committee duty that also holds a malformed duty for a validator the committee has no share for
func StartDutySkipsForeignDuty() tests.SpecTest {

	// KeyShare map with entry only for validator 1
	ksMapFor1Validator := testingutils.KeySetMapForValidators(1)

	duty := testingutils.TestingAttesterDutyForValidators(spec.DataVersionElectra, []int{1})
	// Invalid (slot mismatch, zero pubkey), but validator 2 has no share, so it must be dropped rather than validated
	foreignDuty := &types.ValidatorDuty{
		Type:                    types.BNRoleAttester,
		Slot:                    duty.Slot - 1,
		ValidatorIndex:          2,
		CommitteeIndex:          testingutils.TestingCommitteeIndex,
		CommitteesAtSlot:        testingutils.TestingCommitteesAtSlot,
		CommitteeLength:         testingutils.TestingCommitteeLenght,
		ValidatorCommitteeIndex: testingutils.TestingValidatorCommitteeIndex,
	}
	duty.ValidatorDuties = append([]*types.ValidatorDuty{foreignDuty}, duty.ValidatorDuties...)

	return committee.NewCommitteeSpecTest(
		"start duty skips foreign duty",
		testdoc.CommitteeStartDutySkipsForeignDutyDoc,
		testingutils.BaseCommittee(ksMapFor1Validator),
		[]interface{}{duty},
		"",
		nil,
		nil,
		nil,
		0,
		nil,
	)
}
