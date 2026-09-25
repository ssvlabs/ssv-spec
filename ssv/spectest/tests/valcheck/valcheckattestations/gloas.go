package valcheckattestations

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Gloas (ePBS, SIP #94 §2) twins of the pre-Gloas attestation value-check tests. At a Gloas duty slot
// the committee value is a GloasBeaconVote, so the harness routes these through
// GloasBeaconVoteValueCheckF — which adds the 0/1 attestation-index rule on top of the shared
// checkpoint and slashability checks.

// GloasValid tests a valid GloasBeaconVote.
func GloasValid() tests.SpecTest {
	return valcheck.NewSpecTest(
		"gloas attestation value check valid",
		testdoc.ValCheckGloasAttestationValidDoc,
		types.BeaconTestNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlotGloas,
		testingutils.TestGloasBeaconVoteByts,
		*testingutils.TestGloasBeaconVote.Source,
		*testingutils.TestGloasBeaconVote.Target,
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		0,
		false,
	)
}

// GloasInvalidIndex tests that an attestation data index above 1 is rejected (SIP #94 §2: only
// 0 = payload absent and 1 = payload present are valid).
func GloasInvalidIndex() tests.SpecTest {
	data := testingutils.TestGloasBeaconVote
	data.AttestationDataIndex = 2
	input, err := data.Encode()
	if err != nil {
		panic(err.Error())
	}

	return valcheck.NewSpecTest(
		"gloas attestation value check invalid index",
		testdoc.ValCheckGloasAttestationInvalidIndexDoc,
		types.BeaconTestNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlotGloas,
		input,
		*data.Source,
		*data.Target,
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		types.GloasBeaconVoteInvalidIndexErrorCode,
		false,
	)
}

// GloasSourceHigherThanTarget tests a GloasBeaconVote whose source epoch is not below its target.
func GloasSourceHigherThanTarget() tests.SpecTest {
	data := testingutils.TestGloasBeaconVote
	data.Source = &phase0.Checkpoint{Epoch: 2, Root: testingutils.TestingBlockRoot}
	data.Target = &phase0.Checkpoint{Epoch: 1, Root: testingutils.TestingBlockRoot}
	input, err := data.Encode()
	if err != nil {
		panic(err.Error())
	}

	return valcheck.NewSpecTest(
		"gloas attestation value check source higher than target",
		testdoc.ValCheckGloasAttestationSourceHigherThanTargetDoc,
		types.BeaconTestNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlotGloas,
		input,
		*data.Source,
		*data.Target,
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		types.AttestationSourceNotLessThanTargetErrorCode,
		false,
	)
}

// GloasUnmatchedSourceEpoch tests a GloasBeaconVote whose source epoch differs from the one the
// operator expects from its own view.
func GloasUnmatchedSourceEpoch() tests.SpecTest {
	data := testingutils.TestGloasBeaconVote

	// Different from the vote's source epoch (0), so the checkpoint comparison fails.
	expectedSource := phase0.Checkpoint{Epoch: 1, Root: testingutils.TestingBlockRoot}

	input, err := data.Encode()
	if err != nil {
		panic(err.Error())
	}

	return valcheck.NewSpecTest(
		"gloas attestation value check unmatched source epoch",
		testdoc.ValCheckGloasAttestationUnmatchedSourceEpochDoc,
		types.BeaconTestNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlotGloas,
		input,
		expectedSource,
		*data.Target,
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		types.CheckpointMismatch,
		false,
	)
}

// GloasPreGloasVote tests that a pre-Gloas BeaconVote proposed at a Gloas slot is rejected: the Gloas
// vote is a fixed 120-byte container, so the 112-byte pre-Gloas value fails decoding on length. This
// is the cross-fork decode safety the SIP relies on (SIP #94 §2).
func GloasPreGloasVote() tests.SpecTest {
	return valcheck.NewSpecTest(
		"gloas attestation value check pre-gloas vote",
		testdoc.ValCheckGloasAttestationPreGloasVoteDoc,
		types.BeaconTestNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlotGloas,
		testingutils.TestBeaconVoteByts,
		*testingutils.TestBeaconVote.Source,
		*testingutils.TestBeaconVote.Target,
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		types.DecodeGloasBeaconVoteErrorCode,
		false,
	)
}

// GloasSlashable tests a GloasBeaconVote for a slot already marked slashable for the share.
func GloasSlashable() tests.SpecTest {
	keySet := testingutils.Testing4SharesSet()
	sharePKBytes := keySet.Shares[1].Serialize()

	return valcheck.NewSpecTest(
		"gloas attestation value check slashable",
		testdoc.ValCheckGloasAttestationSlashableDoc,
		types.BeaconTestNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlotGloas,
		testingutils.TestGloasBeaconVoteByts,
		*testingutils.TestGloasBeaconVote.Source,
		*testingutils.TestGloasBeaconVote.Target,
		map[string][]phase0.Slot{
			hex.EncodeToString(sharePKBytes): {
				testingutils.TestingDutySlotGloas,
			},
		},
		[]types.ShareValidatorPK{sharePKBytes},
		types.SlashableAttestationErrorCode,
		false,
	)
}
