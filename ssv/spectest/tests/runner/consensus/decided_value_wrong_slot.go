package consensus

import (
	"crypto/rsa"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DecidedValueWrongSlot tests that a decided value whose duty slot and block slot are shifted together away from
// the running duty's slot is rejected by the value check's running-slot check, so the runner never signs a block
// it is not proposing.
func DecidedValueWrongSlot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := spec.DataVersionDeneb
	duty := testingutils.TestingProposerDutyV(version)

	// A value that is valid on its own terms but for the next slot: duty slot and block slot both shifted.
	wrongCD := testingutils.TestProposerConsensusDataV(version)
	wrongCD.Duty.Slot = duty.Slot + 1
	wrongValue := testingutils.TestProposerConsensusDataMutatedBlockByts(wrongCD, func(slot *phase0.Slot, _ *phase0.ValidatorIndex) {
		*slot = wrongCD.Duty.Slot
	})

	return &tests.MsgProcessingSpecTest{
		Name:          "consensus decided value for wrong slot",
		Documentation: testdoc.ConsensusDecidedValueWrongSlotDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          duty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version))),
			testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
				[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
				[]types.OperatorID{1, 2, 3},
				qbft.Height(duty.Slot),
				testingutils.ProposerMsgID,
				wrongValue,
			),
		},
		// Only the node's own pre-consensus randao share; no post-consensus block signature.
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
		},
		ExpectedErrorCode: types.ProposerDutySlotMismatchErrorCode,
	}
}
