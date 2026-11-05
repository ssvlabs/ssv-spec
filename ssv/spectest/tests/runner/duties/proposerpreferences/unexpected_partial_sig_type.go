package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnexpectedPartialSigType tests that the proposer-preferences runner rejects a pre-consensus partial-signature
// type other than its own two (ProposerPreferencesPartialSig, RequestAuthPartialSig; SIP #94 §5) rather than
// funneling it into the preference path.
func UnexpectedPartialSigType() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	duty := testingutils.TestingProposerPreferencesDuty()

	wrongType := &types.PartialSignatureMessages{
		Type: types.PTCAttesterPartialSig,
		Slot: duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: make([]byte, 96),
				Signer:           1,
				ValidatorIndex:   testingutils.TestingValidatorIndex,
			},
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposer preferences unexpected partial signature type",
		Documentation: testdoc.ProposerPreferencesUnexpectedPartialSigTypeDoc,
		Runner:        testingutils.ProposerPreferencesRunner(ks),
		Duty:          duty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposerPreferences(nil, wrongType)),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		ExpectedErrorCode: types.ProposerPreferencesUnexpectedPartialSigTypeErrorCode,
	}
}
