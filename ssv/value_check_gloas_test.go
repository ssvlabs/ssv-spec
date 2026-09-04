package ssv_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Gloas (ePBS §4): ProposerValueCheckF branches on the fork before any type-layer validation — on a
// Gloas slot it decodes gloas.BeaconBlock directly and never routes through
// ProposerConsensusData.Validate()/GetBlockData() (which have no Gloas arm, since go-eth2-client's
// api.VersionedProposal can't carry Gloas). §2's GloasBeaconVoteValueCheckF is here too. These are
// plain unit tests; the anchor-facing generated vectors are a follow-up (risk-#6 prefix coordination).

// TestProposerConsensusDataValidateErrorsOnGloas documents why the value-check must branch on the
// fork: the types-layer Validate() cannot handle a Gloas value and errors.
func TestProposerConsensusDataValidateErrorsOnGloas(t *testing.T) {
	duty := testingutils.TestingProposerDutyV(gloas.DataVersionGloas)
	block := testingGloasBeaconBlockSSZ(t, duty.Slot)

	// The types-layer Validate() cannot validate a Gloas value (GetBlockData's api.VersionedProposal
	// has no Gloas arm), so it errors — which is exactly why ProposerValueCheckF branches on the fork
	// and never routes a Gloas value through Validate() (SIP #94 §4, branch-and-skip).
	require.Error(t, (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: block}).Validate())
}

// TestProposerValueCheckFGloas asserts the value check decodes a Gloas block: accepting a valid one,
// rejecting an undecodable one (UnmarshalSSZErrorCode), rejecting a slot-mismatched one
// (ProposerBlockSlotMismatchErrorCode), and rejecting a version that disagrees with the duty slot's
// fork in either direction (QBFTValueInvalidErrorCode).
func TestProposerValueCheckFGloas(t *testing.T) {
	km := testingutils.NewTestingKeyManager()
	duty := testingutils.TestingProposerDutyV(gloas.DataVersionGloas)
	valueCheck := ssv.ProposerValueCheckF(km, types.BeaconTestNetwork,
		types.ValidatorPK(testingutils.TestingValidatorPubKey), testingutils.TestingValidatorIndex, nil,
		testingutils.VersionByEpoch)

	valid, err := (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: testingGloasBeaconBlockSSZ(t, duty.Slot)}).Encode()
	require.NoError(t, err)
	require.NoError(t, valueCheck(valid))

	garbage, err := (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: []byte("garbage")}).Encode()
	require.NoError(t, err)
	requireErrorCode(t, valueCheck(garbage), types.UnmarshalSSZErrorCode)

	// a block whose slot does not match the duty slot is rejected
	slotMismatch, err := (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: testingGloasBeaconBlockSSZ(t, duty.Slot+1)}).Encode()
	require.NoError(t, err)
	requireErrorCode(t, valueCheck(slotMismatch), types.ProposerBlockSlotMismatchErrorCode)

	// the leader-stamped Version must agree with the duty slot's fork: a pre-Gloas Version on a Gloas
	// slot is rejected by the explicit guard (the stamp is attacker-controlled) ...
	preGloasVersion, err := (&types.ProposerConsensusData{Duty: *duty, Version: spec.DataVersionElectra, DataSSZ: testingGloasBeaconBlockSSZ(t, duty.Slot)}).Encode()
	require.NoError(t, err)
	requireErrorCode(t, valueCheck(preGloasVersion), types.QBFTValueInvalidErrorCode)

	// ... and a Gloas Version on a pre-Gloas slot takes the pre-Gloas branch, where Validate() rejects
	// it as an unknown version — both mismatch directions fail, matching the node's slot-based check.
	electraDuty := testingutils.TestingProposerDutyV(spec.DataVersionElectra)
	gloasVersionPreGloasSlot, err := (&types.ProposerConsensusData{Duty: *electraDuty, Version: gloas.DataVersionGloas, DataSSZ: testingGloasBeaconBlockSSZ(t, electraDuty.Slot)}).Encode()
	require.NoError(t, err)
	requireErrorCode(t, valueCheck(gloasVersionPreGloasSlot), types.QBFTValueInvalidErrorCode)
}

// testingGloasBeaconBlockSSZ returns the SSZ encoding of the shared minimal Gloas beacon block
// fixture at the given slot.
func testingGloasBeaconBlockSSZ(t *testing.T, slot phase0.Slot) []byte {
	t.Helper()
	byts, err := gloas.TestingBeaconBlock(slot).MarshalSSZ()
	require.NoError(t, err)
	return byts
}

// TestGloasBeaconVoteValueCheckF asserts the §2 Gloas vote value-check: accepts a valid vote,
// rejects index > 1, an undecodable vote, and source >= target.
func TestGloasBeaconVoteValueCheckF(t *testing.T) {
	km := testingutils.NewTestingKeyManager()
	ks := testingutils.Testing4SharesSet()
	sharePubKeys := []types.ShareValidatorPK{ks.Shares[1].GetPublicKey().Serialize()}
	valueCheck := ssv.GloasBeaconVoteValueCheckF(km, testingutils.TestingDutySlotGloas, sharePubKeys,
		testingutils.TestGloasBeaconVote.Source.Epoch, testingutils.TestGloasBeaconVote.Target.Epoch)

	require.NoError(t, valueCheck(testingutils.TestGloasBeaconVoteByts))

	badIndex := testingutils.TestGloasBeaconVote
	badIndex.AttestationDataIndex = 2
	badIndexByts, err := badIndex.Encode()
	require.NoError(t, err)
	requireErrorCode(t, valueCheck(badIndexByts), types.GloasBeaconVoteInvalidIndexErrorCode)

	requireErrorCode(t, valueCheck([]byte("garbage")), types.DecodeGloasBeaconVoteErrorCode)

	badCheckpoints := testingutils.TestGloasBeaconVote
	badCheckpoints.Source = &phase0.Checkpoint{Epoch: 2, Root: testingutils.TestingBlockRoot}
	badCheckpoints.Target = &phase0.Checkpoint{Epoch: 1, Root: testingutils.TestingBlockRoot}
	badCheckpointsByts, err := badCheckpoints.Encode()
	require.NoError(t, err)
	requireErrorCode(t, valueCheck(badCheckpointsByts), types.AttestationSourceNotLessThanTargetErrorCode)
}

// requireErrorCode asserts err is a *types.Error carrying the given code.
func requireErrorCode(t *testing.T, err error, code int) {
	t.Helper()
	var e *types.Error
	require.ErrorAs(t, err, &e)
	require.Equal(t, code, e.Code)
}

// TestGloasBeaconVoteValueCheckF_CrossIndexEquivocation pins the one semantic rule Gloas adds beyond a
// shape change (SIP #94 §2): the decided attestation index goes into the slashability data, replacing
// the pre-Gloas math.MaxUint64 sentinel. Two votes over the same (source, target, slot) that differ
// only in that index are then a double vote; under the sentinel they would build identical attestation
// data and the equivocation would go undetected — a slashing-safety failure rather than a liveness one.
//
// The key manager records the first attestation data it is asked about per slot and reports a later,
// differing one as slashable.
//
// Only the last assertion discriminates: reinstating the sentinel makes it fail, because both votes then
// build identical attestation data and the equivocation disappears. The one before it — that repeating
// the same vote is not equivocation — holds under either implementation, so it guards against
// over-eager slashing rather than pinning this rule; do not read it as protecting the index behaviour.
func TestGloasBeaconVoteValueCheckF_CrossIndexEquivocation(t *testing.T) {
	km := testingutils.NewTestingKeyManagerRecordingAttestations()
	ks := testingutils.Testing4SharesSet()
	sharePubKeys := []types.ShareValidatorPK{ks.Shares[1].GetPublicKey().Serialize()}

	payloadPresent := testingutils.TestGloasBeaconVote // AttestationDataIndex = 1
	require.EqualValues(t, 1, payloadPresent.AttestationDataIndex)
	payloadAbsent := payloadPresent
	payloadAbsent.AttestationDataIndex = 0

	valueCheck := ssv.GloasBeaconVoteValueCheckF(km, testingutils.TestingDutySlotGloas, sharePubKeys,
		payloadPresent.Source.Epoch, payloadPresent.Target.Epoch)

	presentByts, err := payloadPresent.Encode()
	require.NoError(t, err)
	absentByts, err := payloadAbsent.Encode()
	require.NoError(t, err)

	// The operator votes payload-present.
	require.NoError(t, valueCheck(presentByts))

	// Re-voting the very same value is not equivocation.
	require.NoError(t, valueCheck(presentByts))

	// Voting payload-absent for the same (source, target, slot) is: only the index differs.
	requireErrorCode(t, valueCheck(absentByts), types.SlashableAttestationErrorCode)
}
