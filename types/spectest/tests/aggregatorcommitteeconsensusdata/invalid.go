package aggregatorcommitteeconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidAggregatorValidationNoValidators tests an invalid consensus data with no validators
func InvalidAggregatorValidationNoValidators() *AggregatorCommitteeConsensusDataTest {

	cd := types.AggregatorCommitteeConsensusData{
		Version: gloas.DataVersionGloas,
	}
	return NewAggregatorCommitteeConsensusDataTest(
		"invalid aggregator data with no validators",
		testdoc.AggregatorCommitteeConsensusDataTestInvalidNoValidatorsDoc,
		cd,
		types.AggCommConsensusDataNoValidatorErrorCode,
	)
}

// InvalidAggregatorValidationCommitteeIndexesLength tests an invalid consensus data with wrong CommitteeIndexes length
func InvalidAggregatorValidationCommitteeIndexesLength() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(gloas.DataVersionGloas)

	cd.AggregatorsCommitteeIndexes = append(cd.AggregatorsCommitteeIndexes, cd.AggregatorsCommitteeIndexes...)

	return NewAggregatorCommitteeConsensusDataTest(
		"invalid aggregator data with wrong committee indexes length",
		testdoc.AggregatorCommitteeConsensusDataTestInvalidCommitteeIndexLenDoc,
		*cd,
		types.AggCommAggCommIdxCntMismatchErrorCode,
	)
}

// InvalidAggregatorValidationDuplicateCommitteeIndex tests an invalid consensus data with duplicated CommitteeIndex
func InvalidAggregatorValidationDuplicateCommitteeIndex() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(gloas.DataVersionGloas)

	cd.AggregatorsCommitteeIndexes = append(cd.AggregatorsCommitteeIndexes, cd.AggregatorsCommitteeIndexes[0])
	cd.AggregatedAttestations = append(cd.AggregatedAttestations, cd.AggregatedAttestations[0])

	return NewAggregatorCommitteeConsensusDataTest(
		"invalid aggregator data with duplicated committee index",
		testdoc.AggregatorCommitteeConsensusDataTestDuplicatedCommitteeIndexDoc,
		*cd,
		types.AggCommDuplicatedCommIdxErrorCode,
	)
}

// InvalidAggregatorValidationMissingCommitteeIndex tests an invalid consensus data in which an aggregator's committee index is missing from the existing CommitteeIndex set
func InvalidAggregatorValidationMissingCommitteeIndex() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(gloas.DataVersionGloas)

	maxCommIndex := cd.AggregatorsCommitteeIndexes[0]
	for _, idx := range cd.AggregatorsCommitteeIndexes {
		if idx > maxCommIndex {
			maxCommIndex = idx
		}
	}
	cd.Aggregators[0].CommitteeIndex = maxCommIndex + 1

	return NewAggregatorCommitteeConsensusDataTest(
		"invalid aggregator data with missing committee index",
		testdoc.AggregatorCommitteeConsensusDataTestMissingCommitteeIndexDoc,
		*cd,
		types.AggCommCommIdxMismatchErrorCode,
	)
}

// InvalidAggregatorValidationUnusedCommitteeIndex tests an invalid consensus data in which a committee index is left unused
func InvalidAggregatorValidationUnusedCommitteeIndex() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(gloas.DataVersionGloas)

	maxCommIndex := cd.AggregatorsCommitteeIndexes[0]
	for _, idx := range cd.AggregatorsCommitteeIndexes {
		if idx > maxCommIndex {
			maxCommIndex = idx
		}
	}
	cd.Aggregators[0].CommitteeIndex = maxCommIndex + 1
	cd.AggregatorsCommitteeIndexes = append(cd.AggregatorsCommitteeIndexes, maxCommIndex+1)
	cd.AggregatedAttestations = append(cd.AggregatedAttestations, cd.AggregatedAttestations[0])

	return NewAggregatorCommitteeConsensusDataTest(
		"invalid aggregator data with unused committee index",
		testdoc.AggregatorCommitteeConsensusDataTestUnusedCommitteeIndexDoc,
		*cd,
		types.AggCommUnusedCommIdxErrorCode,
	)
}

// InvalidAggregatorValidationAttestationDecoding tests an invalid consensus data in which an attestation fails to decode
func InvalidAggregatorValidationAttestationDecoding() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(gloas.DataVersionGloas)

	cd.AggregatedAttestations[0] = []byte{0x01, 0x02, 0x03} // invalid attestation bytes

	return NewAggregatorCommitteeConsensusDataTest(
		"invalid aggregator data with attestation decoding error",
		testdoc.AggregatorCommitteeConsensusDataTestInvalidAttestationDecodingDoc,
		*cd,
		types.AggCommAttestationDecodingErrorCode,
	)
}

// InvalidAggregatorValidationNonGloasVersion tests an invalid consensus data whose Version is not Gloas
func InvalidAggregatorValidationNonGloasVersion() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(gloas.DataVersionGloas)

	cd.Version = gloas.DataVersionGloas - 1

	return NewAggregatorCommitteeConsensusDataTest(
		"invalid aggregator data with non-gloas version",
		testdoc.AggregatorCommitteeConsensusDataTestInvalidNonGloasVersionDoc,
		*cd,
		types.UnknownVersionErrorCode,
	)
}

// InvalidSyncCommitteeContributionSubnet tests an invalid consensus data with a duplicated subnet
func InvalidSyncCommitteeContributionDuplicatedSubnet() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestSyncCommitteeContributionConsensusDataF()

	cd.SyncCommitteeContributions = append(cd.SyncCommitteeContributions, cd.SyncCommitteeContributions[0])

	return NewAggregatorCommitteeConsensusDataTest(
		"invalid sync committee contribution with duplicated subnet",
		testdoc.AggregatorCommitteeConsensusDataTestInvalidDuplicatedSubnetDoc,
		*cd,
		types.AggCommSCCSubnetDuplicateErrorCode,
	)
}

// InvalidSyncCommitteeContributionSubnet tests an invalid consensus data with a missing subnet
func InvalidSyncCommitteeContributionMissingSubnet() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestSyncCommitteeContributionConsensusDataF()

	cd.Contributors[0].CommitteeIndex = 100

	return NewAggregatorCommitteeConsensusDataTest(
		"invalid sync committee contribution with missing subnet",
		testdoc.AggregatorCommitteeConsensusDataTestMissingSubnetDoc,
		*cd,
		types.AggCommSubnetNotInSCSubnetsErrorCode,
	)
}

// InvalidSyncCommitteeContributionUnusedSubnet tests an invalid consensus data with an unused subnet
func InvalidSyncCommitteeContributionUnusedSubnet() *AggregatorCommitteeConsensusDataTest {

	cd := testingutils.TestSyncCommitteeContributionConsensusDataF()

	unused := cd.SyncCommitteeContributions[0]
	unused.SubcommitteeIndex = 100
	cd.SyncCommitteeContributions = append(cd.SyncCommitteeContributions, unused)

	return NewAggregatorCommitteeConsensusDataTest(
		"invalid sync committee contribution with unused subnet",
		testdoc.AggregatorCommitteeConsensusDataTestUnusedSubnetDoc,
		*cd,
		types.AggCommUnusedSubnetErrorCode,
	)
}
