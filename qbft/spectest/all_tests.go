package spectest

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/messages"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	messages.CommitDataEncoding(),
	messages.DecidedMsgEncoding(),
	messages.MsgNilIdentifier(),
	messages.MsgNonZeroIdentifier(),
	messages.MsgTypeUnknown(),
	messages.PrepareDataEncoding(),
	messages.ProposeDataEncoding(),
	messages.MsgDataNil(),
	messages.MsgDataNonZero(),
	//messages.RoundChangeDataEncoding(),
	messages.SignedMsgSigTooShort(),
	messages.SignedMsgSigTooLong(),
	messages.SignedMsgNoSigners(),
	messages.GetRoot(),
	messages.SignedMessageEncoding(),

	//tests.HappyFlow(),
	//tests.SevenOperators(),
	//tests.TenOperators(),
	//tests.ThirteenOperators(),
	//
	//commit.HappyFlow(),
	//commit.MultiSignerWithOverlap(),
	//commit.MultiSignerNoOverlap(),
	//commit.Decided(),
	//commit.NoPrevAcceptedProposal(),
	//commit.WrongHeight(),
	//commit.WrongRound(),
	//commit.ImparsableCommitData(),
	//commit.WrongCommitData(),
	//commit.WrongSignature(),
}
