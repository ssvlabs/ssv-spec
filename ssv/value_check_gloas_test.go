package ssv_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	bitfield "github.com/prysmaticlabs/go-bitfield"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Gloas (ePBS §4, design "Option A"): ProposerConsensusData.Validate treats the Gloas block as
// opaque (keeping the types package free of the gloas import), while ProposerValueCheckF owns the
// gloas.BeaconBlock decode. These live as a plain unit test because the value-check's Gloas decode
// is not yet expressible as a generated spectest vector (the ssv/spectest suite is blocked, and the
// proposerconsensusdata suite only exercises the now-opaque Validate); the anchor-facing vector is
// a follow-up once that is unblocked and the fixture prefix is coordinated with anchor.

// TestProposerConsensusDataValidateGloasOpaque asserts Validate is opaque for a Gloas value: it
// accepts any DataSSZ (block decode/validation is the value check's job) but still checks the duty.
func TestProposerConsensusDataValidateGloasOpaque(t *testing.T) {
	duty := testingutils.TestingProposerDutyV(gloas.DataVersionGloas)
	block := testingGloasBeaconBlockSSZ(t, duty.Slot)

	require.NoError(t, (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: block}).Validate())
	// Opaque: even an undecodable block passes Validate (rejected later, in the value check).
	require.NoError(t, (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: []byte("garbage")}).Validate())

	// The duty-type check still runs for Gloas.
	wrongDuty := *duty
	wrongDuty.Type = types.BNRoleAttester
	require.Error(t, (&types.ProposerConsensusData{Duty: wrongDuty, Version: gloas.DataVersionGloas, DataSSZ: block}).Validate())
}

// TestProposerValueCheckFGloas asserts the value check decodes a Gloas block: accepting a valid one,
// rejecting an undecodable one (UnmarshalSSZErrorCode), and rejecting a slot-mismatched one
// (ProposerBlockSlotMismatchErrorCode).
func TestProposerValueCheckFGloas(t *testing.T) {
	km := testingutils.NewTestingKeyManager()
	duty := testingutils.TestingProposerDutyV(gloas.DataVersionGloas)
	valueCheck := ssv.ProposerValueCheckF(km, types.BeaconTestNetwork,
		types.ValidatorPK(testingutils.TestingValidatorPubKey), testingutils.TestingValidatorIndex, nil)

	valid, err := (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: testingGloasBeaconBlockSSZ(t, duty.Slot)}).Encode()
	require.NoError(t, err)
	require.NoError(t, valueCheck(valid))

	garbage, err := (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: []byte("garbage")}).Encode()
	require.NoError(t, err)
	var decodeErr *types.Error
	require.ErrorAs(t, valueCheck(garbage), &decodeErr)
	require.Equal(t, types.UnmarshalSSZErrorCode, decodeErr.Code)

	// a block whose slot does not match the duty slot is rejected
	slotMismatch, err := (&types.ProposerConsensusData{Duty: *duty, Version: gloas.DataVersionGloas, DataSSZ: testingGloasBeaconBlockSSZ(t, duty.Slot+1)}).Encode()
	require.NoError(t, err)
	var slotErr *types.Error
	require.ErrorAs(t, valueCheck(slotMismatch), &slotErr)
	require.Equal(t, types.ProposerBlockSlotMismatchErrorCode, slotErr.Code)
}

// testingGloasBeaconBlockSSZ builds a minimal valid Gloas beacon block (mirrors the
// types/gloas round-trip fixture) and returns its SSZ encoding.
func testingGloasBeaconBlockSSZ(t *testing.T, slot phase0.Slot) []byte {
	t.Helper()
	byts, err := (&gloas.BeaconBlock{
		Slot:          slot,
		ProposerIndex: testingutils.TestingValidatorIndex,
		Body: &gloas.BeaconBlockBody{
			ETH1Data:      &phase0.ETH1Data{BlockHash: make([]byte, 32)},
			SyncAggregate: &altair.SyncAggregate{SyncCommitteeBits: bitfield.NewBitvector512()},
			SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{
				BuilderIndex: gloas.BuilderIndexSelfBuild,
			}},
			ParentExecutionRequests: &gloas.ExecutionRequests{},
		},
	}).MarshalSSZ()
	require.NoError(t, err)
	return byts
}
