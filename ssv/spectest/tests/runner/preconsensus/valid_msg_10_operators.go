package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidMessage10Operators tests a valid PartialSignatureMessages with multi PartialSignatureMessages (10 operators)
func ValidMessage10Operators() tests.SpecTest {
	ks := testingutils.Testing10SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus valid msg 10 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgSyncCommitteeContribution(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommitteeContribution(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgSyncCommitteeContribution(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgSyncCommitteeContribution(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[5], 5)),
					testingutils.SSVMsgSyncCommitteeContribution(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[6], 6)),
				},
				PostDutyRunnerStateRoot: "b56aa8d90c9d249bde18de5ce59dfdcef7c681d9c27f776b125074365a35d618",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgAggregator(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgAggregator(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgAggregator(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgAggregator(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgAggregator(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[5], 5)),
					testingutils.SSVMsgAggregator(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[6], 6)),
				},
				PostDutyRunnerStateRoot: "f92ab2b0c921bc10053ed7e2323598192766310910d00b9a32dc2a81a5175474",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], 5, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[6], 6, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "5d129bfe44158b4abc0eae740f7f00e5023bbfdd467de80a62b1a984ff34bce4",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], 5, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(6, ks.NetworkKeys[6], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[6], 6, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "6eba81cc6ec2322051ac347dc3d15c9501f12de0fdd5918fc825edc57412923e",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
		},
	}
}
