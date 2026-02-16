package valcheckaggcommittee

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoValidator tests no validators
func NoValidator() tests.SpecTest {

	cd := types.AggregatorCommitteeConsensusData{}
	cdBytes, err := cd.Encode()
	if err != nil {
		panic(err.Error())
	}

	return valcheck.NewSpecTest(
		"aggcommittee value check no validators",
		testdoc.ValCheckAggCommitteeNoValidatorsDoc,
		types.PraterNetwork,
		types.RoleAggregatorCommittee,
		testingutils.TestingDutySlot,
		cdBytes,
		phase0.Checkpoint{},
		phase0.Checkpoint{},
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		types.AggCommConsensusDataNoValidatorErrorCode,
		true,
	)
}
