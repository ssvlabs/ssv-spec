package valcheckaggcommittee

import (
	"github.com/OffchainLabs/go-bitfield"
	eth2gloas "github.com/attestantio/go-eth2-client/spec/gloas"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// electraMaxAggregationBits is the bound on an Electra attestation's aggregation bits:
// MAX_VALIDATORS_PER_COMMITTEE * MAX_COMMITTEES_PER_SLOT.
const electraMaxAggregationBits = 131072

// LongAggregationBits tests that the aggregate is decoded as the Gloas attestation (SIP #94 §2), whose
// aggregation bits are a progressive bitlist with no bound of its own: an aggregate one bit longer than an
// Electra attestation may be is accepted. A value check that still decodes the Electra container rejects it.
func LongAggregationBits() tests.SpecTest {
	version := gloas.DataVersionGloas
	consensusData := testingutils.TestAggregatorCommitteeConsensusDataForDuty(testingutils.TestingAggregatorCommitteeDutyOnlyAggregator(version), version, nil)

	aggregate := &eth2gloas.Attestation{}
	if err := aggregate.UnmarshalSSZ(consensusData.AggregatedAttestations[0]); err != nil {
		panic(err.Error())
	}
	aggregate.AggregationBits = bitfield.NewBitlist(electraMaxAggregationBits + 1)
	encoded, err := aggregate.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	consensusData.AggregatedAttestations[0] = encoded

	input, err := consensusData.Encode()
	if err != nil {
		panic(err.Error())
	}

	return valcheck.NewMultiSpecTest(
		"aggcommittee value check long aggregation bits",
		testdoc.ValCheckAggCommitteeLongAggregationBitsDoc,
		[]*valcheck.SpecTest{
			{
				Name:       "one bit over the electra bound",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      input,
			},
		},
	)
}
