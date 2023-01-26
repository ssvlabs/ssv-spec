package spectest

import (
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/messages"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{

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

	messages.CreateVCBC(),
	messages.CreateABA(),
	messages.CreateFillGap(),
	messages.CreateFiller(),
	messages.CreateABAInit(),
	messages.CreateABAAux(),
	messages.CreateABAConf(),
	messages.CreateABAFinish(),
	messages.CreateVCBCBroadcast(),
	messages.CreateVCBCSend(),
	messages.CreateVCBCReady(),
	messages.CreateVCBCFinal(),
	messages.CreateVCBCRequest(),
	messages.CreateVCBCAnswer(),
}
