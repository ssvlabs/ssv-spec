package valcheckproposer

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasBlockListLimits covers the list limits of a proposed Gloas block (SIP #94 §4; see
// gloas.BeaconBlock.ValidateListLimits). The value check rejects a block the beacon chain would reject for a
// list over its limit, so that the round changes instead of the cluster signing a block that loses the slot at
// publication.
//
// One case holds every list at its limit and is accepted. Each other case has one list a single element over
// its limit, in an otherwise empty block, and is rejected.
func GloasBlockListLimits() tests.SpecTest {
	duty := testingutils.TestingProposerDutyV(gloas.DataVersionGloas)

	value := func(lists gloas.TestingBlockLists) []byte {
		dataSSZ, err := (&gloas.GloasProposalData{
			Block:       gloas.TestingBeaconBlockWithLists(duty.Slot, lists),
			PayloadRoot: testingutils.TestingGloasPayloadRoot,
		}).MarshalSSZ()
		if err != nil {
			panic(err.Error())
		}
		byts, err := (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: dataSSZ}).Encode()
		if err != nil {
			panic(err.Error())
		}
		return byts
	}
	over := func(name string, lists gloas.TestingBlockLists) *valcheck.SpecTest {
		return &valcheck.SpecTest{
			Name:              name,
			Network:           types.BeaconTestNetwork,
			RunnerRole:        types.RoleProposer,
			Input:             value(lists),
			ExpectedErrorCode: types.QBFTValueInvalidErrorCode,
		}
	}

	// The counts are written out, and not taken from the gloas limit constants, so that the vectors pin them.
	atLimits := gloas.TestingBlockLists{
		ProposerSlashings:      16,
		AttesterSlashings:      1,
		Attestations:           8,
		Deposits:               0,
		VoluntaryExits:         16,
		BLSToExecutionChanges:  16,
		PayloadAttestations:    4,
		WithdrawalRequests:     16,
		ConsolidationRequests:  2,
		BuilderDepositRequests: 64,
		BuilderExitRequests:    16,
		BlobKZGCommitments:     4096,
	}

	return valcheck.NewMultiSpecTest(
		"gloas block list limits",
		testdoc.ValCheckProposerGloasBlockListLimitsDoc,
		[]*valcheck.SpecTest{
			{
				Name:       "every list at its limit",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleProposer,
				Input:      value(atLimits),
			},
			over("proposer slashings over the limit", gloas.TestingBlockLists{ProposerSlashings: 17}),
			over("attester slashings over the limit", gloas.TestingBlockLists{AttesterSlashings: 2}),
			over("attestations over the limit", gloas.TestingBlockLists{Attestations: 9}),
			// The body carries no deposits at Gloas: they reach the chain as execution requests.
			over("a deposit", gloas.TestingBlockLists{Deposits: 1}),
			over("voluntary exits over the limit", gloas.TestingBlockLists{VoluntaryExits: 17}),
			over("bls to execution changes over the limit", gloas.TestingBlockLists{BLSToExecutionChanges: 17}),
			over("payload attestations over the limit", gloas.TestingBlockLists{PayloadAttestations: 5}),
			over("withdrawal requests over the limit", gloas.TestingBlockLists{WithdrawalRequests: 17}),
			over("consolidation requests over the limit", gloas.TestingBlockLists{ConsolidationRequests: 3}),
			over("builder deposit requests over the limit", gloas.TestingBlockLists{BuilderDepositRequests: 65}),
			over("builder exit requests over the limit", gloas.TestingBlockLists{BuilderExitRequests: 17}),
			over("blob kzg commitments over the limit", gloas.TestingBlockLists{BlobKZGCommitments: 4097}),
		},
	)
}
