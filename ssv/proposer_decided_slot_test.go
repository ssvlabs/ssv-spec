package ssv_test

import (
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TestProposerRejectsDecidedValueForWrongSlot covers the SIP #94 §4 duty-slot binding. QBFT binds the
// decided message to the running instance height, but nothing binds the decided value's own Duty.Slot to
// the running duty, so without the ProcessConsensus guard a leader could get a value for another slot
// decided and make operators sign a block for a duty they are not running. The runner runs slot Y with an
// instance at height Y; a decided (quorum-commit) message carrying a valid value for slot Y+1 is injected
// at height Y, so only the value's slot is wrong and the guard must reject it.
func TestProposerRejectsDecidedValueForWrongSlot(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	runningDuty := testingutils.TestingProposerDutyV(version)
	runningSlot := runningDuty.Slot
	wrongSlot := runningSlot + 1

	// A value that is valid on its own terms but for a different slot: an external-build block (payload_root
	// zero passes the value check's self-build pin), for slot Y+1, with a matching duty slot.
	wrongDuty := testingutils.TestingProposerDutyV(version)
	wrongDuty.Slot = wrongSlot
	proposalData := &gloas.GloasProposalData{Block: gloas.TestingBeaconBlockExternalBuild(wrongSlot)}
	dataSSZ, err := proposalData.MarshalSSZ()
	require.NoError(t, err)
	wrongValue := &types.ProposerConsensusData{Duty: *wrongDuty, Version: version, DataSSZ: dataSSZ}
	wrongFullData, err := wrongValue.Encode()
	require.NoError(t, err)

	msgID := types.NewValidatorMsgID(testingutils.TestingSSVDomainType, types.ValidatorPK(testingutils.TestingValidatorPubKey), types.RoleProposer)

	// Start the duty and complete the randao round so the runner produces its own value and starts a
	// consensus instance at height == runningSlot (the first ks.Threshold messages are the randao partials).
	v := testingutils.BaseValidator(ks)
	require.NoError(t, v.StartDuty(runningDuty))
	randaoMsgs := testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer)[:ks.Threshold]
	for _, msg := range randaoMsgs {
		require.NoError(t, v.ProcessMessage(msg))
	}

	// Inject a decided value for the wrong slot at the running instance's height.
	decided := testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
		[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
		[]types.OperatorID{1, 2, 3},
		qbft.Height(runningSlot),
		msgID[:],
		wrongFullData,
	)
	err = v.ProcessMessage(decided)
	require.Error(t, err)
	var typedErr *types.Error
	require.ErrorAs(t, err, &typedErr)
	require.Equal(t, types.ProposerDecidedSlotMismatchErrorCode, typedErr.Code)
}
