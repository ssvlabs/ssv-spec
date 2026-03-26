package preconsensus

import (
	"crypto/sha256"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostDecided tests a msg received post consensus decided (and post receiving a quorum for pre consensus)
func PostDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	numValidators := 1
	validators := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	decideRunner := func(r ssv.Runner, duty *types.ValidatorDuty, decidedValue *types.ProposerConsensusData, preMsgs []*types.PartialSignatureMessages) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		for _, msg := range preMsgs {
			err := r.ProcessPreConsensus(msg)
			if err != nil {
				panic(err.Error())
			}
		}
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.CommitteeMember,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight,
			r.GetBaseRunner().QBFTController.OperatorSigner)
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().State.DecidedValue = testingutils.EncodeConsensusDataTest(decidedValue)
		r.GetBaseRunner().QBFTController.StoredInstances[0] = r.GetBaseRunner().State.RunningInstance
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
		return r
	}

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus post decided",
		testdoc.PreConsensusPostDecidedDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	// Aggregator committee duty

	// SC Duty
	scDuty := testingutils.TestingAggregatorCommitteeDuty([]int{}, validators, spec.DataVersionAltair)
	height := qbft.Height(scDuty.Slot)
	msgID := testingutils.AggregatorCommitteeMsgID(ks)
	scConsensusData := testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(scDuty, spec.DataVersionAltair)

	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name:   "sync committee aggregator selection proof",
		Runner: testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
		Duty:   scDuty,
		Messages: []*types.SignedSSVMessage{
			// Normal pre-consensus messages
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(scDuty, ksMap, 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(scDuty, ksMap, 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(scDuty, ksMap, 3))),

			// Consensus messages
			testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), msgID, scConsensusData, height),
			testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(scConsensusData)),
			testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(scConsensusData)),
			testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(scConsensusData)),
			testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(scConsensusData)),
			testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(scConsensusData)),
			testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(scConsensusData)),

			// Extra pre-consensus message after consensus decided should be ignored
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(scDuty, ksMap, 4))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusAggregatorCommitteeMsgForDuty(scDuty, ksMap, 1),
			testingutils.PostConsensusAggregatorCommitteeMsgForDuty(scDuty, ksMap, 1, spec.DataVersionAltair),
		},
	})
	for _, version := range testingutils.SupportedAggregatorVersions {
		// Aggregator duty
		aggDuty := testingutils.TestingAggregatorCommitteeDuty(validators, []int{}, version)
		aggHeight := qbft.Height(aggDuty.Slot)
		aggConsensusData := testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(aggDuty, version)

		// Mixed duty
		mixedDuty := testingutils.TestingAggregatorCommitteeDutyMixed(version)
		mixedConsensusData := testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(mixedDuty, version)
		mixedHeight := qbft.Height(mixedDuty.Slot)

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:   fmt.Sprintf("aggregator selection proof (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
				Duty:   aggDuty,
				Messages: []*types.SignedSSVMessage{
					// Normal pre-consensus messages
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 3))),

					// Consensus messages
					testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), msgID, aggConsensusData, aggHeight),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, aggHeight, msgID, sha256.Sum256(aggConsensusData)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, aggHeight, msgID, sha256.Sum256(aggConsensusData)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, aggHeight, msgID, sha256.Sum256(aggConsensusData)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, aggHeight, msgID, sha256.Sum256(aggConsensusData)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, aggHeight, msgID, sha256.Sum256(aggConsensusData)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, aggHeight, msgID, sha256.Sum256(aggConsensusData)),

					// Extra pre-consensus message after consensus decided should be ignored
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 4))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1),
					testingutils.PostConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1, version),
				},
			},
			{
				Name:   fmt.Sprintf("aggregator committee duty (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
				Duty:   mixedDuty,
				Messages: []*types.SignedSSVMessage{
					// Normal pre-consensus messages
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 3))),

					// Consensus messages
					testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), msgID, mixedConsensusData, mixedHeight),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, mixedHeight, msgID, sha256.Sum256(mixedConsensusData)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, mixedHeight, msgID, sha256.Sum256(mixedConsensusData)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, mixedHeight, msgID, sha256.Sum256(mixedConsensusData)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, mixedHeight, msgID, sha256.Sum256(mixedConsensusData)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, mixedHeight, msgID, sha256.Sum256(mixedConsensusData)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, mixedHeight, msgID, sha256.Sum256(mixedConsensusData)),

					// Extra pre-consensus message after consensus decided should be ignored
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 4))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1),
					testingutils.PostConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1, version),
				},
			},
		}...)
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("randao (%s)", version.String()),
			Runner: decideRunner(
				testingutils.ProposerRunner(ks),
				testingutils.TestingProposerDutyV(version),
				testingutils.TestProposerConsensusDataV(version),
				[]*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version),
					testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, version),
					testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, version),
				},
			),
			Duty: testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], ks.Shares[4], 4, 4, version))),
			},
			DontStartDuty: true,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("randao blinded block (%s)", version.String()),
			Runner: decideRunner(
				testingutils.ProposerBlindedBlockRunner(ks),
				testingutils.TestingProposerDutyV(version),
				testingutils.TestProposerBlindedBlockConsensusDataV(version),
				[]*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version),
					testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, version),
					testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, version),
				},
			),
			Duty: testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], ks.Shares[4], 4, 4, version))),
			},
			DontStartDuty: true,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
