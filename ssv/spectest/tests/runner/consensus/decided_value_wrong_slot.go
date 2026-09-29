package consensus

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DecidedValueWrongSlot tests that a decided value for a slot other than the running duty's is rejected by the
// value check's running-slot bind (SIP #94 §4), so the runner never signs a block it is not proposing.
func DecidedValueWrongSlot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas
	duty := testingutils.TestingProposerDutyV(version)
	wrongSlot := duty.Slot + 1

	// A value that is valid on its own terms but for the next slot: an external-bid block (zero payload_root
	// passes the self-build rule) for that slot, with a matching duty slot.
	wrongDuty := testingutils.TestingProposerDutyV(version)
	wrongDuty.Slot = wrongSlot
	dataSSZ, err := (&gloas.GloasProposalData{Block: gloas.TestingBeaconBlockExternalBuild(wrongSlot)}).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	wrongValue, err := (&types.ProposerConsensusData{Duty: *wrongDuty, Version: version, DataSSZ: dataSSZ}).Encode()
	if err != nil {
		panic(err.Error())
	}

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
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
		},
		ExpectedErrorCode: types.ProposerDutySlotMismatchErrorCode,
	}
}
