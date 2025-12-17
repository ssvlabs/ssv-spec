package aggregatorcommitteeconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidAggregatorValidationAttestationLength tests an invalid consensus data with wrong Attestation length
func InvalidAggregatorValidationAttestationLength() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0)

	cd.Attestations = append(cd.Attestations, cd.Attestations...) // duplicate to make length invalid

	return NewValidatorConsensusDataTest(
		"invalid aggregator data with wrong attestation length",
		testdoc.AggregatorCommitteeConsensusDataTestInvalidAttestationLenDoc,
		*cd,
		types.AggCommAggAttCntMismatchErrorCode,
	)
}

// InvalidAggregatorValidationCommitteeIndexesLength tests an invalid consensus data with wrong CommitteeIndexes length
func InvalidAggregatorValidationCommitteeIndexesLength() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0)

	cd.AggregatorsCommitteeIndexes = append(cd.AggregatorsCommitteeIndexes, cd.AggregatorsCommitteeIndexes...)

	return NewValidatorConsensusDataTest(
		"invalid aggregator data with wrong committee indexes length",
		testdoc.AggregatorCommitteeConsensusDataTestInvalidCommitteeIndexLenDoc,
		*cd,
		types.AggCommAggCommIdxCntMismatchErrorCode,
	)
}

// InvalidSyncCommitteeContributionLength tests an invalid consensus data with invalid sync committee contrib length
func InvalidSyncCommitteeContributionLength() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestSyncCommitteeContributionConsensusDataF()

	cd.SyncCommitteeContributions = append(cd.SyncCommitteeContributions, cd.SyncCommitteeContributions...)

	return NewValidatorConsensusDataTest(
		"invalid sync committee contribution with wrong contribution length",
		testdoc.AggregatorCommitteeConsensusDataTestInvalidContributionLenDoc,
		*cd,
		types.AggCommContributorsContributionsCntMismatchErrorCode,
	)
}

// InvalidSyncCommitteeContributionSubnet tests an invalid consensus data with a subnet mismatch
func InvalidSyncCommitteeContributionSubnet() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestSyncCommitteeContributionConsensusDataF()

	cd.SyncCommitteeContributions[0].SubcommitteeIndex = 100

	return NewValidatorConsensusDataTest(
		"invalid sync committee contribution with subnet mismatch",
		testdoc.AggregatorCommitteeConsensusDataTestInvalidSubnetMismatchDoc,
		*cd,
		types.AggCommSubnetNotInSCSubnetsErrorCode,
	)
}
