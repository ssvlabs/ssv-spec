package valcheckattestations

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SameSlotIndex tests the same-slot rule (SIP #94 §2): a vote for a block proposed in the attestation's own slot
// must carry index 0. An operator whose own view establishes that the voted block is of the duty's slot rejects
// index 1 for it. It accepts index 0 for it, index 1 for any other block, and, with no such view, index 1 for
// any block.
func SameSlotIndex() tests.SpecTest {
	slot := phase0.Slot(testingutils.TestingDutySlot)

	payloadPresent := testingutils.TestBeaconVote // AttestationDataIndex = 1
	payloadAbsent := payloadPresent
	payloadAbsent.AttestationDataIndex = 0

	encode := func(bv types.BeaconVote) []byte {
		byts, err := bv.Encode()
		if err != nil {
			panic(err.Error())
		}
		return byts
	}

	votedBlock := []phase0.Root{payloadPresent.BlockRoot}
	otherBlock := []phase0.Root{{0x0b, 0x10, 0xc8}}

	test := func(name string, vote types.BeaconVote, sameSlotBlockRoots []phase0.Root, expectedErrorCode int) *valcheck.SpecTest {
		return &valcheck.SpecTest{
			Name:               name,
			Network:            types.BeaconTestNetwork,
			RunnerRole:         types.RoleCommittee,
			DutySlot:           slot,
			Input:              encode(vote),
			ExpectedSource:     *vote.Source,
			ExpectedTarget:     *vote.Target,
			SameSlotBlockRoots: sameSlotBlockRoots,
			ExpectedErrorCode:  expectedErrorCode,
		}
	}

	return valcheck.NewMultiSpecTest(
		"attestation value check same-slot index",
		testdoc.ValCheckAttestationSameSlotIndexDoc,
		[]*valcheck.SpecTest{
			test("index 1 for the block of the slot", payloadPresent, votedBlock, types.BeaconVoteSameSlotIndexErrorCode),
			test("index 0 for the block of the slot", payloadAbsent, votedBlock, 0),
			test("index 1 for another block", payloadPresent, otherBlock, 0),
			test("index 1 without a view", payloadPresent, nil, 0),
		},
	)
}
