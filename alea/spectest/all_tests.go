package spectest

import (
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/abaaux"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/abaconf"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/abafinish"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/abainit"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/filler"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/fillgap"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/messages"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/proposal"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/vcbcanswer"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/vcbcfinal"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/vcbcready"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/vcbcrequest"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/vcbcsend"
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
	messages.MsgDataNil(),
	messages.MsgDataNonZero(),
	messages.SignedMsgSigTooShort(),
	messages.SignedMsgSigTooLong(),
	messages.SignedMsgNoSigners(),
	messages.SignedMsgDuplicateSigners(),
	messages.SignedMsgMultiSigners(),
	messages.ProposalDataInvalid(),
	messages.SignedMessageSigner0(),

	messages.ProposeDataEncoding(),
	messages.SignedMessageEncoding(),
	messages.GetRoot(),

	messages.CreateProposal(),
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
	proposal.EmptyData(),

	vcbcsend.MultiSigner(),
	vcbcsend.UnknownSigner(),
	vcbcsend.WrongHeight(),
	vcbcsend.WrongSignature(),
	vcbcsend.EmptyData(),

	vcbcready.EmptyHash(),
	vcbcready.MultiSigner(),
	vcbcready.UnknownSigner(),
	vcbcready.WrongHeight(),
	vcbcready.WrongSignature(),

	vcbcfinal.MultiSigner(),
	vcbcfinal.UnknownSigner(),
	vcbcfinal.WrongHeight(),
	vcbcfinal.WrongSignature(),
	vcbcfinal.EmptyAggregatedMsgBytes(),
	vcbcfinal.EmptyHash(),

	vcbcrequest.MultiSigner(),
	vcbcrequest.UnknownSigner(),
	vcbcrequest.WrongHeight(),
	vcbcrequest.WrongSignature(),

	vcbcanswer.MultiSigner(),
	vcbcanswer.UnknownSigner(),
	vcbcanswer.WrongHeight(),
	vcbcanswer.WrongSignature(),

	abainit.MultiSigner(),
	abainit.UnknownSigner(),
	abainit.WrongHeight(),
	abainit.WrongSignature(),

	abaaux.MultiSigner(),
	abaaux.UnknownSigner(),
	abaaux.WrongHeight(),
	abaaux.WrongSignature(),

	abaconf.MultiSigner(),
	abaconf.UnknownSigner(),
	abaconf.WrongHeight(),
	abaconf.WrongSignature(),

	abafinish.MultiSigner(),
	abafinish.UnknownSigner(),
	abafinish.WrongHeight(),
	abafinish.WrongSignature(),

	filler.MultiSigner(),
	filler.UnknownSigner(),
	filler.WrongHeight(),
	filler.WrongSignature(),

	fillgap.MultiSigner(),
	fillgap.UnknownSigner(),
	fillgap.WrongHeight(),
	fillgap.WrongSignature(),
}
