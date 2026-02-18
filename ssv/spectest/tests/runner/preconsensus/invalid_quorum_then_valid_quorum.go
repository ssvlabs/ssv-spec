package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidQuorumThenValid tests a runner receiving an invalid message forming an invalid quorum, then receiving a valid message forming a valid quorum, terminating successfully
func InvalidQuorumThenValidQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedErrorCode := types.ReconstructSignatureErrorCode

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus invalid quorum then valid quorum",
		testdoc.PreConsensusInvalidQuorumThenValidQuorumDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationWrongBeaconSigMsg(ks.Shares[1], 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[4], 4))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedValidatorRegistration(ks)),
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitWrongBeaconSigMsg(ks.Shares[1], 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[2], 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[3], 3))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[4], 4))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedVoluntaryExit(ks)),
				},
				ExpectedErrorCode: expectedErrorCode,
			},
		},
		ks,
	)

	// Aggregator committee duty
	sccDuty := testingutils.TestingSyncCommitteeContributionDuty
	sccSlot := sccDuty.Slot
	msgID := testingutils.TestingAggregatorCommitteeMsgID
	sccConsensusData := testingutils.TestAggregatorCommitteeConsensusDataForDuty(sccDuty, spec.DataVersionPhase0, nil)
	sccCDBytes, err := sccConsensusData.Encode()
	if err != nil {
		panic(err)
	}
	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name:   "sync committee aggregator selection proof",
		Runner: testingutils.AggregatorCommitteeRunner(ks),
		Duty:   sccDuty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofWrongBeaconSigMsg(ks.Shares[1], ks.Shares[1], 1, 1, sccSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[2], ks.Shares[2], 2, 2, sccSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[3], ks.Shares[3], 3, 3, sccSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[4], ks.Shares[4], 4, 4, sccSlot))),
			// Proposal msg to be accepted
			testingutils.SignedSSVMessageWithSignerAndFullData(1, ks.OperatorKeys[1], testingutils.SSVMsgSyncCommitteeContribution(testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], 1, msgID[:], sccCDBytes, qbft.Height(sccSlot)), nil), sccCDBytes),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusContributionProofMsgWithSlot(ks.Shares[1], ks.Shares[1], 1, 1, sccSlot), // broadcasts when starting a new duty
		},
		ExpectedErrorCode: expectedErrorCode,
	})
	for _, version := range testingutils.SupportedAggregatorVersions {
		// Agg
		aggDuty := testingutils.TestingAggregatorDuty(version)
		aggSlot := aggDuty.Slot
		aggConsensusData := testingutils.TestAggregatorCommitteeConsensusDataForDuty(aggDuty, version, nil)
		aggCDBytes, err := aggConsensusData.Encode()
		if err != nil {
			panic(err)
		}
		// Mixed
		mixedDuty := testingutils.TestingAggregatorCommitteeDutyMixed(version)
		mixedSlot := mixedDuty.Slot
		mixedConsensusData := testingutils.TestAggregatorCommitteeConsensusDataForDuty(mixedDuty, version, nil)
		mixedCDBytes, err := mixedConsensusData.Encode()
		if err != nil {
			panic(err)
		}

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:   fmt.Sprintf("aggregator selection proof (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofWrongBeaconSigMsg(ks.Shares[1], ks.Shares[1], 1, 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], ks.Shares[4], 4, 4, version))),
					// Proposal
					testingutils.SignedSSVMessageWithSignerAndFullData(1, ks.OperatorKeys[1], testingutils.SSVMsgAggregator(testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], 1, msgID[:], aggCDBytes, qbft.Height(aggSlot)), nil), aggCDBytes),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingAggregatorCommitteeDutyMixed(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsgWrongBeaconSig(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[3], 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[4], 4, version))),
					// Proposal
					testingutils.SignedSSVMessageWithSignerAndFullData(1, ks.OperatorKeys[1], testingutils.SSVMsgAggregatorCommittee(ks, testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], 1, msgID[:], mixedCDBytes, qbft.Height(mixedSlot)), nil), mixedCDBytes),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusAggregatorCommitteeMixedMsg(ks.Shares[1], 1, version), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: expectedErrorCode,
			},
		}...)
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		duty := testingutils.TestingProposerDutyV(version)
		slot := duty.Slot
		cd := testingutils.TestProposerConsensusDataV(version)
		cdBytes, err := cd.Encode()
		if err != nil {
			panic(err)
		}
		msgID := testingutils.ProposerMsgID
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("randao (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoWrongBeaconSigMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[4], 4, version))),
				// Proposal
				testingutils.SignedSSVMessageWithSignerAndFullData(1, ks.OperatorKeys[1], testingutils.SSVMsgProposer(testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], 1, msgID[:], cdBytes, qbft.Height(slot)), nil), cdBytes),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
			ExpectedErrorCode: expectedErrorCode,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		duty := testingutils.TestingProposerDutyV(version)
		slot := duty.Slot
		cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
		cdBytes, err := cd.Encode()
		if err != nil {
			panic(err)
		}
		msgID := testingutils.ProposerMsgID
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("randao blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoWrongBeaconSigMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[4], 4, version))),
				// Proposal
				testingutils.SignedSSVMessageWithSignerAndFullData(1, ks.OperatorKeys[1], testingutils.SSVMsgProposer(testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], 1, msgID[:], cdBytes, qbft.Height(slot)), nil), cdBytes),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
			ExpectedErrorCode: expectedErrorCode,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
