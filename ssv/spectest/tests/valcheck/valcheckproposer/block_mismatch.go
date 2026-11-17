package valcheckproposer

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ProposerBlockSlotMismatch tests a proposer value whose block slot differs from the duty slot
func ProposerBlockSlotMismatch() tests.SpecTest {
	return proposerBlockMismatch(
		"proposer block slot mismatch",
		testdoc.ValCheckProposerBlockSlotMismatchDoc,
		func(slot *phase0.Slot, _ *phase0.ValidatorIndex) { *slot++ },
		types.ProposerBlockSlotMismatchErrorCode,
	)
}

// ProposerBlockProposerIndexMismatch tests a proposer value whose block proposer index differs from the duty
// validator index
func ProposerBlockProposerIndexMismatch() tests.SpecTest {
	return proposerBlockMismatch(
		"proposer block proposer index mismatch",
		testdoc.ValCheckProposerBlockProposerIndexMismatchDoc,
		func(_ *phase0.Slot, index *phase0.ValidatorIndex) { *index++ },
		types.ProposerBlockProposerIndexMismatchErrorCode,
	)
}

func proposerBlockMismatch(
	name, doc string,
	mutate func(*phase0.Slot, *phase0.ValidatorIndex),
	expectedErrCode int,
) tests.SpecTest {
	specTests := make([]*valcheck.SpecTest, 0)
	forEachProposerValue(func(name string, cd *types.ProposerConsensusData) {
		specTests = append(specTests, &valcheck.SpecTest{
			Name:              name,
			Network:           types.BeaconTestNetwork,
			RunnerRole:        types.RoleProposer,
			Input:             testingutils.TestProposerConsensusDataMutatedBlockByts(cd, mutate),
			ExpectedErrorCode: expectedErrCode,
		})
	})
	return valcheck.NewMultiSpecTest(name, doc, specTests)
}

// RunningDutySlotMismatch tests a proposer value whose duty slot and block slot are shifted together away from
// the running duty's slot. The block still matches the value's duty, so only the running-slot check rejects it.
func RunningDutySlotMismatch() tests.SpecTest {
	specTests := make([]*valcheck.SpecTest, 0)
	forEachProposerValue(func(name string, cd *types.ProposerConsensusData) {
		runningSlot := cd.Duty.Slot
		specTests = append(specTests, &valcheck.SpecTest{
			Name:       name + " running slot",
			Network:    types.BeaconTestNetwork,
			RunnerRole: types.RoleProposer,
			DutySlot:   runningSlot,
			Input:      testingutils.TestProposerConsensusDataMutatedBlockByts(cd, func(*phase0.Slot, *phase0.ValidatorIndex) {}),
		})
		for _, delta := range []int64{1, 40, -1000} {
			shifted := *cd
			shifted.Duty.Slot = phase0.Slot(int64(runningSlot) + delta)
			specTests = append(specTests, &valcheck.SpecTest{
				Name:       fmt.Sprintf("%s slot %+d", name, delta),
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleProposer,
				DutySlot:   runningSlot,
				Input: testingutils.TestProposerConsensusDataMutatedBlockByts(&shifted, func(slot *phase0.Slot, _ *phase0.ValidatorIndex) {
					*slot = shifted.Duty.Slot
				}),
				ExpectedErrorCode: types.ProposerDutySlotMismatchErrorCode,
			})
		}
	})
	return valcheck.NewMultiSpecTest(
		"proposer running duty slot mismatch",
		testdoc.ValCheckProposerRunningDutySlotMismatchDoc,
		specTests,
	)
}

// forEachProposerValue calls f with a fresh valid proposer value for every supported block version, regular and
// blinded, named after the role, e.g. "proposer (capella)" and "proposer blinded block (capella)".
func forEachProposerValue(f func(name string, cd *types.ProposerConsensusData)) {
	for _, version := range testingutils.SupportedBlockVersions {
		f(fmt.Sprintf("proposer (%s)", version), testingutils.TestProposerConsensusDataV(version))
		f(fmt.Sprintf("proposer blinded block (%s)", version), testingutils.TestProposerBlindedBlockConsensusDataV(version))
	}
}
