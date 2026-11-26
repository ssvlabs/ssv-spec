package proposerconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasBlockValidation tests a valid consensus data with a Gloas block (SIP #94 §4)
func GloasBlockValidation() *ProposerConsensusDataTest {
	return NewProposerConsensusDataTest(
		"valid gloas block",
		testdoc.ProposerConsensusDataTestGloasBlockDoc,
		*testingutils.TestProposerConsensusDataV(gloas.DataVersionGloas),
		0,
	)
}

// InvalidGloasBlockValidation tests an invalid consensus data with empty Gloas proposal data
func InvalidGloasBlockValidation() *ProposerConsensusDataTest {
	version := gloas.DataVersionGloas

	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewProposerConsensusDataTest(
		"invalid gloas block",
		testdoc.ProposerConsensusDataTestInvalidGloasBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}

// NonGloasVersionValidation tests an invalid consensus data whose Version is not Gloas, carrying valid Gloas
// proposal data (SIP #94 §4 pins the version)
func NonGloasVersionValidation() *ProposerConsensusDataTest {
	cd := testingutils.TestProposerConsensusDataV(gloas.DataVersionGloas)
	cd.Version = gloas.DataVersionGloas - 1

	return NewProposerConsensusDataTest(
		"non-gloas version",
		testdoc.ProposerConsensusDataTestNonGloasVersionDoc,
		*cd,
		types.UnknownBlockVersionErrorCode,
	)
}

// GloasBlockListOverLimitValidation tests an invalid consensus data whose block holds one attestation more than
// a Gloas block may. The list is progressive, so the block decodes, and the limit is checked on its own
// (SIP #94 §4).
func GloasBlockListOverLimitValidation() *ProposerConsensusDataTest {
	version := gloas.DataVersionGloas
	duty := testingutils.TestingProposerDutyV(version)

	dataSSZ, err := (&gloas.GloasProposalData{
		Block:       gloas.TestingBeaconBlockWithLists(duty.Slot, gloas.TestingBlockLists{Attestations: 9}), // one over MAX_ATTESTATIONS_ELECTRA
		PayloadRoot: testingutils.TestingGloasPayloadRoot,
	}).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return NewProposerConsensusDataTest(
		"gloas block list over limit",
		testdoc.ProposerConsensusDataTestGloasBlockListOverLimitDoc,
		types.ProposerConsensusData{Duty: *duty, Version: version, DataSSZ: dataSSZ},
		types.ProposerBlockListLimitErrorCode,
	)
}
