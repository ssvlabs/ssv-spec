package spectest

import (
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/messages"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/proposal"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{

	tests.HappyFlow(),
	tests.SevenOperators(),
	tests.TenOperators(),
	tests.ThirteenOperators(),

	messages.MsgNilIdentifier(),
	messages.MsgNonZeroIdentifier(),
	messages.MsgTypeUnknown(),
	messages.ProposeDataEncoding(),
	messages.MsgDataNil(),
	messages.MsgDataNonZero(),
	messages.SignedMsgSigTooShort(),
	messages.SignedMsgSigTooLong(),
	messages.SignedMsgNoSigners(),
	messages.SignedMsgDuplicateSigners(),
	messages.SignedMsgMultiSigners(),
	messages.GetRoot(),
	messages.SignedMessageEncoding(),
	messages.CreateProposal(),
	messages.ProposalDataInvalid(),
	messages.SignedMessageSigner0(),

	messages.CreateFillGap(),
	messages.CreateFiller(),
	messages.CreateABAInit(),
	messages.CreateABAAux(),
	messages.CreateABAConf(),
	messages.CreateABAFinish(),
	messages.CreateVCBCSend(),
	messages.CreateVCBCReady(),
	messages.CreateVCBCFinal(),
	messages.CreateVCBCRequest(),
	messages.CreateVCBCAnswer(),

	proposal.MultiSigner(),
	proposal.PostDecided(),
	proposal.UnknownSigner(),
	proposal.WrongHeight(),
	proposal.WrongSignature(),
}
