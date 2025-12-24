package valcheckaggcommittee

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidAttLen tests for the invalid attestation length case
func InvalidAttLen() tests.SpecTest {

	cd := testingutils.TestAggregatorConsensusData(spec.DataVersionElectra)
	cd.Attestations = append(cd.Attestations, cd.Attestations[0])
	cdBytes, err := cd.Encode()
	if err != nil {
		panic(err.Error())
	}

	return valcheck.NewSpecTest(
		"aggcommittee value check invalid attestation length",
		testdoc.ValCheckAggCommitteeInvalidAttestationLenDoc,
		types.PraterNetwork,
		types.RoleAggregatorCommittee,
		testingutils.TestingDutySlot,
		cdBytes,
		phase0.Checkpoint{},
		phase0.Checkpoint{},
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		types.AggCommAggAttCntMismatchErrorCode,
		true,
	)
}
