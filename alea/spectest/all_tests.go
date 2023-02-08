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
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests/flow"
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
	vcbcsend.Batch(),
	vcbcsend.TwoBatches(),
	vcbcsend.Receive(),
	vcbcsend.ReceiveMultiple(),
	vcbcsend.WrongAuthor(),
	vcbcsend.Duplicate(),
	vcbcsend.DuplicatePriority(),

	vcbcready.EmptyHash(),
	vcbcready.MultiSigner(),
	vcbcready.UnknownSigner(),
	vcbcready.WrongHeight(),
	vcbcready.WrongSignature(),
	vcbcready.Receive(),
	vcbcready.ReceiveQuorum(),
	vcbcready.WrongHash(),
	vcbcready.UnexpectedAuthor(),
	vcbcready.Duplicate(),

	vcbcfinal.MultiSigner(),
	vcbcfinal.UnknownSigner(),
	vcbcfinal.WrongHeight(),
	vcbcfinal.WrongSignature(),
	vcbcfinal.EmptyAggregatedMsgBytes(),
	vcbcfinal.EmptyHash(),
	vcbcfinal.Receive(),
	vcbcfinal.ReceiveRequest(),
	vcbcfinal.Duplicate(),
	vcbcfinal.WrongProof(),

	vcbcrequest.MultiSigner(),
	vcbcrequest.UnknownSigner(),
	vcbcrequest.WrongHeight(),
	vcbcrequest.WrongSignature(),
	vcbcrequest.Receive(),
	vcbcrequest.WrongAuthor(),
	vcbcrequest.UnknownPriority(),

	vcbcanswer.MultiSigner(),
	vcbcanswer.UnknownSigner(),
	vcbcanswer.WrongHeight(),
	vcbcanswer.WrongSignature(),
	vcbcanswer.EmptyData(),
	vcbcanswer.WrongData(),
	vcbcanswer.WrongPriority(),
	vcbcanswer.WrongAnswer(),

	abainit.MultiSigner(),
	abainit.UnknownSigner(),
	abainit.WrongHeight(),
	abainit.WrongSignature(),
	abainit.InvalidVote(),
	abainit.Duplicate(),
	abainit.Receive(),
	abainit.ReceiveQuorum(),
	abainit.ReceiveTwoQuorum(),
	abainit.AbaStart(),

	abaaux.MultiSigner(),
	abaaux.UnknownSigner(),
	abaaux.WrongHeight(),
	abaaux.WrongSignature(),
	abaaux.InvalidVote(),
	abaaux.Duplicate(),
	abaaux.Receive(),
	abaaux.ReceiveQuorum(),
	abaaux.ReceiveTwoQuorum(),
	abaaux.ReceiveNoQuorum(),

	abaconf.MultiSigner(),
	abaconf.UnknownSigner(),
	abaconf.WrongHeight(),
	abaconf.WrongSignature(),
	abaconf.InvalidVote(),
	abaconf.Receive(),
	abaconf.ReceiveQuorum(),
	abaconf.Receive2ValuesQuorum(),
	abaconf.ReceiveNoQuorum(),

	abafinish.MultiSigner(),
	abafinish.UnknownSigner(),
	abafinish.WrongHeight(),
	abafinish.WrongSignature(),
	abafinish.InvalidVote(),
	abafinish.Receive(),
	abafinish.ReceiveQuorum(),
	abafinish.ReceiveNoQuorum(),
	abafinish.Duplicate(),

	filler.MultiSigner(),
	filler.UnknownSigner(),
	filler.WrongHeight(),
	filler.WrongSignature(),
	filler.WrongData(),
	filler.WrongPriority(),
	filler.EmptyData(),
	filler.WrongAnswer(),

	fillgap.MultiSigner(),
	fillgap.UnknownSigner(),
	fillgap.WrongHeight(),
	fillgap.WrongSignature(),
	fillgap.Receive(),
	fillgap.UnknownPriority(),
	fillgap.WrongAuthor(),

	flow.Flow7Op(),
}
