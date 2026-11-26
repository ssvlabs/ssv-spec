package ptcattester

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ptcMsg operator-signs a peer's PTC partial-signature container.
func ptcMsg(ks *testingutils.TestKeySet, msg *types.PartialSignatureMessages) *types.SignedSSVMessage {
	return testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPTCAttester(nil, msg))
}

// RetainedPartialsAppliedAtDutyStart tests the partials that arrive before this operator's observation is
// frozen (SIP #94 §3): a peer signs as soon as its own observation is final, which can be before this
// operator's duty starts. Three peers' partials arrive first and are retained; starting the duty freezes the
// observation and replays them, reaching quorum and submitting the payload attestation message.
func RetainedPartialsAppliedAtDutyStart() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MsgProcessingSpecTest{
		Name:          "ptc attestation retained partials applied at duty start",
		Documentation: testdoc.PTCAttesterRetainedPartialsAppliedAtDutyStartDoc,
		Runner:        testingutils.PTCAttesterRunner(ks),
		Duty:          testingutils.TestingPTCAttesterDuty(),
		PreDutyMessages: []*types.SignedSSVMessage{
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[2], 2)),
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[3], 3)),
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[4], 4)),
		},
		Messages: []*types.SignedSSVMessage{}, // the duty start replays the retained partials
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusPTCMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedPayloadAttestationMessage(ks)),
		},
	}
}

// RetainedDivergingPartialDropped tests that retained partials (SIP #94 §3) are checked against the
// observation once it is frozen: the first, from a peer that observed otherwise, fails the expected-root check
// and is skipped without failing the duty, while the two matching ones behind it still apply. This operator's
// own partial then completes the quorum.
func RetainedDivergingPartialDropped() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MsgProcessingSpecTest{
		Name:          "ptc attestation retained diverging partial dropped",
		Documentation: testdoc.PTCAttesterRetainedDivergingPartialDroppedDoc,
		Runner:        testingutils.PTCAttesterRunner(ks),
		Duty:          testingutils.TestingPTCAttesterDuty(),
		PreDutyMessages: []*types.SignedSSVMessage{
			ptcMsg(ks, testingutils.PreConsensusPTCWrongRootMsg(ks.Shares[2], 2)),
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[3], 3)),
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[4], 4)),
		},
		Messages: []*types.SignedSSVMessage{
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[1], 1)),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusPTCMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		BeaconBroadcastedRoots: []string{
			// retained op3 + op4 + this operator's op1 share = quorum; op2's diverging share was skipped.
			testingutils.GetSSZRootNoError(testingutils.TestingSignedPayloadAttestationMessage(ks)),
		},
	}
}

// RetainedPartialsOnePerSigner tests that a signer has one retained partial per slot (SIP #94 §3/§7): its
// first is kept and a later one is not. Operator 2 sends a diverging partial and then a matching one before the
// duty starts. Only the diverging one is replayed, and it is skipped, so the matching partials of operators 3
// and 4 stay one short of quorum and nothing is submitted.
func RetainedPartialsOnePerSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MsgProcessingSpecTest{
		Name:          "ptc attestation retained partials one per signer",
		Documentation: testdoc.PTCAttesterRetainedPartialsOnePerSignerDoc,
		Runner:        testingutils.PTCAttesterRunner(ks),
		Duty:          testingutils.TestingPTCAttesterDuty(),
		PreDutyMessages: []*types.SignedSSVMessage{
			ptcMsg(ks, testingutils.PreConsensusPTCWrongRootMsg(ks.Shares[2], 2)),
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[2], 2)),
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[3], 3)),
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[4], 4)),
		},
		Messages: []*types.SignedSSVMessage{},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusPTCMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
	}
}

// RetainedPartialReconstructFailureReturned tests that replaying retained partials (SIP #94 §3) reports a
// quorum it completes but cannot reconstruct, as the live path does: operator 4's retained share has a bad
// signature, so the replayed quorum fails to reconstruct, the duty start returns that error, and the bad share
// is dropped. This operator's own share then re-crosses quorum and the message is submitted.
func RetainedPartialReconstructFailureReturned() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MsgProcessingSpecTest{
		Name:          "ptc attestation retained partial reconstruct failure returned",
		Documentation: testdoc.PTCAttesterRetainedPartialReconstructFailureReturnedDoc,
		Runner:        testingutils.PTCAttesterRunner(ks),
		Duty:          testingutils.TestingPTCAttesterDuty(),
		PreDutyMessages: []*types.SignedSSVMessage{
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[2], 2)),
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[3], 3)),
			ptcMsg(ks, testingutils.PreConsensusPTCWrongBeaconSigMsg(ks.Shares[4], 4)),
		},
		Messages: []*types.SignedSSVMessage{
			ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[1], 1)),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusPTCMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		BeaconBroadcastedRoots: []string{
			// op2 + op3 + this operator's op1 share, once op4's bad share is dropped.
			testingutils.GetSSZRootNoError(testingutils.TestingSignedPayloadAttestationMessage(ks)),
		},
		ExpectedErrorCode: types.ReconstructSignatureErrorCode,
	}
}

// RetainedPartialsBySlot tests which duty start uses, keeps or drops a slot's retained partials (SIP #94 §3):
// an abstaining duty reconstructs nothing from them; a duty for an earlier slot leaves them for their own; a
// duty for a later slot drops them, as their own can no longer start.
func RetainedPartialsBySlot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return tests.NewMultiMsgProcessingSpecTest(
		"ptc attestation retained partials by slot",
		testdoc.PTCAttesterRetainedPartialsBySlotDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				// A quorum of peers observed a block where this operator's beacon node reports none.
				Name:   "abstaining duty",
				Runner: testingutils.PTCAttesterRunner(ks),
				Duty:   testingutils.TestingPTCAttesterAbstainDuty(),
				PreDutyMessages: []*types.SignedSSVMessage{
					ptcMsg(ks, testingutils.PreConsensusPTCAbstainSlotMsg(ks.Shares[2], 2)),
					ptcMsg(ks, testingutils.PreConsensusPTCAbstainSlotMsg(ks.Shares[3], 3)),
					ptcMsg(ks, testingutils.PreConsensusPTCAbstainSlotMsg(ks.Shares[4], 4)),
				},
				Messages:       []*types.SignedSSVMessage{},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
			{
				// The retained partials are for the next epoch's duty; the post-state still holds them.
				Name:   "duty for an earlier slot",
				Runner: testingutils.PTCAttesterRunner(ks),
				Duty:   testingutils.TestingPTCAttesterDuty(),
				PreDutyMessages: []*types.SignedSSVMessage{
					ptcMsg(ks, testingutils.PreConsensusPTCNextEpochMsg(ks.Shares[2], 2)),
					ptcMsg(ks, testingutils.PreConsensusPTCNextEpochMsg(ks.Shares[3], 3)),
					ptcMsg(ks, testingutils.PreConsensusPTCNextEpochMsg(ks.Shares[4], 4)),
				},
				Messages: []*types.SignedSSVMessage{},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusPTCMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				// The retained partials are for a slot the started duty is already past; the post-state holds none.
				Name:   "duty for a later slot",
				Runner: testingutils.PTCAttesterRunner(ks),
				Duty:   testingutils.TestingPTCAttesterNextEpochDuty(),
				PreDutyMessages: []*types.SignedSSVMessage{
					ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[2], 2)),
					ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[3], 3)),
					ptcMsg(ks, testingutils.PreConsensusPTCMsg(ks.Shares[4], 4)),
				},
				Messages: []*types.SignedSSVMessage{},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusPTCNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
		ks,
	)
}

// MalformedPartialsNotRetained tests that a malformed packet arriving before the duty starts is rejected and
// not retained (SIP #94 §3/§7): one with more than the single partial of a PTC duty, and one for a validator
// this runner does not run. The post-state pins that nothing was retained.
func MalformedPartialsNotRetained() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	unknownValidator := testingutils.PreConsensusPTCMsg(ks.Shares[2], 2)
	unknownValidator.Messages[0].ValidatorIndex = testingutils.TestingWrongValidatorIndex

	return tests.NewMultiMsgProcessingSpecTest(
		"ptc attestation malformed partials not retained",
		testdoc.PTCAttesterMalformedPartialsNotRetainedDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "two partials",
				Runner: testingutils.PTCAttesterRunner(ks),
				Duty:   testingutils.TestingPTCAttesterDuty(),
				Messages: []*types.SignedSSVMessage{
					ptcMsg(ks, testingutils.PreConsensusPTCTooManyRootsMsg(ks.Shares[2], 2)),
				},
				OutputMessages:    []*types.PartialSignatureMessages{},
				DontStartDuty:     true,
				ExpectedErrorCode: types.WrongRootsCountErrorCode,
			},
			{
				Name:   "partial for unknown validator",
				Runner: testingutils.PTCAttesterRunner(ks),
				Duty:   testingutils.TestingPTCAttesterDuty(),
				Messages: []*types.SignedSSVMessage{
					ptcMsg(ks, unknownValidator),
				},
				OutputMessages:    []*types.PartialSignatureMessages{},
				DontStartDuty:     true,
				ExpectedErrorCode: types.UnknownValidatorIndexErrorCode,
			},
		},
		ks,
	)
}
