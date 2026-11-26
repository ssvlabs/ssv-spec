package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerExternalBid tests the external-build path end to end (SIP #94 §4): an external bid win returns a
// bare block and no envelope, so the runner proposes a value with a zero payload_root. Once that value is
// decided, the post-consensus packet carries the block entry alone, the block is submitted, and no reveal is
// published: the envelope is the builder's to reveal.
func GloasProposerExternalBid() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	consensusData := testingutils.TestProposerExternalBidConsensusDataV(version)
	value, err := consensusData.Encode()
	if err != nil {
		panic(err.Error())
	}
	proposal := &gloas.GloasProposalData{Block: gloas.TestingBeaconBlockExternalBuild(testingutils.TestingDutySlotV(version))}

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer external bid",
		Documentation: testdoc.GloasProposerExternalBidDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages: append(
			testingutils.SSVDecidingMsgsV(consensusData, ks, types.RoleProposer), // pre-consensus + consensus
			[]*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgFor(ks.Shares[1], 1, proposal))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgFor(ks.Shares[2], 2, proposal))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgFor(ks.Shares[3], 3, proposal))),
			}...,
		),
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			testingutils.PostConsensusProposerMsgFor(ks.Shares[1], 1, proposal), // the block entry alone
		},
		QBFTProposals: [][]byte{value},
		// The block only: an external bid has no envelope for the cluster to reveal.
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockOf(ks, proposal)),
		},
		BeaconNode: &tests.BeaconNodeBehaviour{ProducesExternalBid: true},
	}
}
