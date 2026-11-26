package runner

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DomainEpoch tests the epoch the Gloas duties sign for (SIP #94 §3/§5/§6): the epoch of the duty's slot. For
// the proposer-preferences duty that slot is the proposal slot, however early the duty is emitted, and the
// dependent root it pins is that epoch's proposer-duties root. The request-auth domain is chain-independent.
//
// Each duty runs its happy flow at TestingDutySlotEpochDependent, in an epoch for which the testing beacon node
// has a fork version and a dependent root unlike any other epoch's: a signing root computed for another epoch
// matches none of the peers' partials, and reconstructs nothing that is expected here.
func DomainEpoch() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	slot := phase0.Slot(testingutils.TestingDutySlotEpochDependent)

	quorum := []types.OperatorID{1, 2, 3}

	// Proposer: the randao reveal, the block root and the §6 envelope root.
	proposerDuty := testingutils.TestingProposerDutyForSlot(slot)
	proposal := testingutils.TestingGloasProposalData(slot)
	consensusData := &types.ProposerConsensusData{
		Duty:    *proposerDuty,
		Version: gloas.DataVersionGloas,
		DataSSZ: testingutils.TestingGloasProposalDataBytes(slot),
	}
	var proposerMsgs []*types.SignedSSVMessage
	for _, opID := range quorum {
		proposerMsgs = append(proposerMsgs, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgForSlot(ks.Shares[opID], opID, slot))))
	}
	proposerMsgs = append(proposerMsgs, testingutils.SSVDecidingMsgsForHeight(consensusData, testingutils.ProposerMsgID, qbft.Height(slot), ks)...)
	for _, opID := range quorum {
		proposerMsgs = append(proposerMsgs, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgFor(ks.Shares[opID], opID, proposal))))
	}

	// PTC attestation: the payload attestation data.
	var ptcMsgs []*types.SignedSSVMessage
	for _, opID := range quorum {
		ptcMsgs = append(ptcMsgs, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPTCAttester(nil, testingutils.PreConsensusPTCMsgForSlot(ks.Shares[opID], opID, slot))))
	}

	// Proposer preferences: the preference, with its dependent root, and the two request auths.
	authData := [][]byte{testingutils.TestingBuilderEntries[0].AuthData(), testingutils.TestingBuilderEntries[1].AuthData()}
	var preferencesMsgs []*types.SignedSSVMessage
	for _, opID := range quorum {
		preferencesMsgs = append(preferencesMsgs, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusProposerPreferencesMsgForSlot(ks.Shares[opID], opID, slot))))
	}
	for _, opID := range quorum {
		preferencesMsgs = append(preferencesMsgs, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, testingutils.PreConsensusBuilderRequestAuthMsgForSlot(ks.Shares[opID], opID, authData, slot))))
	}

	return tests.NewMultiMsgProcessingSpecTest(
		"domain epoch",
		testdoc.DomainEpochDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:     "proposer",
				Runner:   testingutils.ProposerRunner(ks),
				Duty:     proposerDuty,
				Messages: proposerMsgs,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgForSlot(ks.Shares[1], 1, slot),
					testingutils.PostConsensusProposerMsgFor(ks.Shares[1], 1, proposal),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockOf(ks, proposal)),
					testingutils.TestingSignedEnvelopeRootOf(ks, proposal),
				},
				OrderedBeaconBroadcastedRoots: true,
			},
			{
				Name:     "ptc attestation",
				Runner:   testingutils.PTCAttesterRunner(ks),
				Duty:     testingutils.TestingPTCAttesterDutyForSlot(slot),
				Messages: ptcMsgs,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusPTCMsgForSlot(ks.Shares[1], 1, slot),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedPayloadAttestationMessageForSlot(ks, slot)),
				},
			},
			{
				Name:     "proposer preferences",
				Runner:   testingutils.ProposerPreferencesRunnerWithBuilderEntries(ks),
				Duty:     testingutils.TestingProposerPreferencesDutyForSlot(slot),
				Messages: preferencesMsgs,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusProposerPreferencesMsgForSlot(ks.Shares[1], 1, slot),
					testingutils.PreConsensusBuilderRequestAuthMsgForSlot(ks.Shares[1], 1, authData, slot),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedProposerPreferences(ks, slot)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, authData[0], slot)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedBuilderRequestAuth(ks, authData[1], slot)),
				},
			},
		},
		ks,
	)
}
