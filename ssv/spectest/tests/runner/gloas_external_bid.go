package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerExternalBid tests the external-build produce path (SIP #94 §4): an external bid win returns a bare
// block and no envelope, so the runner proposes a value with a zero payload_root.
func GloasProposerExternalBid() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	proposal, err := testingutils.TestProposerExternalBidConsensusDataV(version).Encode()
	if err != nil {
		panic(err.Error())
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer external bid",
		Documentation: testdoc.GloasProposerExternalBidDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
		},
		QBFTProposals: [][]byte{proposal},
		BeaconNode:    &tests.BeaconNodeBehaviour{ProducesExternalBid: true},
	}
}
